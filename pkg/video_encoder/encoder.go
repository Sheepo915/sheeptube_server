package video_encoder

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
)

type EncodingJob struct {
	RunID      string
	VideoID    string
	BucketName string
	ObjectName string
}

type Encoder interface {
	Encode(ctx context.Context, job *EncodingJob) error
}

type VideoEncoder struct {
	tempPath    string
	minioClient *minio.Client
	opt         *VideoEncoderOption
}

type VideoEncoderOption struct {
	MaxWorker int
	*RetryPolicy
}

type RetryPolicy struct {
	InitialInterval    time.Duration
	BackoffCoefficient float64
	MaxInterval        time.Duration
	MaxAttempts        int
}

func NewEncoder(mc *minio.Client, option *VideoEncoderOption) *VideoEncoder {
	opt := &VideoEncoderOption{
		MaxWorker: 5,
		RetryPolicy: &RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaxInterval:        100 * time.Second,
			MaxAttempts:        0,
		},
	}

	if option != nil {
		opt = option
	}

	return &VideoEncoder{
		tempPath:    "/temp",
		minioClient: mc,
		opt:         opt,
	}
}

func (e *VideoEncoder) SetTempPath(path string) {
	e.tempPath = path
}

func (e *VideoEncoder) Encode(ctx context.Context, job *EncodingJob) error {
	attempt := 0
	interval := e.opt.InitialInterval

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := e.EncodeVideo(ctx, job)
		if err == nil {
			return nil // success!
		}

		attempt++
		if e.opt.MaxAttempts > 0 && attempt >= e.opt.MaxAttempts {
			return fmt.Errorf("exceeded max attempts: %w", err)
		}

		fmt.Printf("Encoding failed (attempt %d): %v. Retrying in %v...\n", attempt, err, interval)

		// Sleep before retrying
		time.Sleep(interval)

		// Increase interval with backoff (capped)
		interval = time.Duration(float64(interval) * float64(e.opt.BackoffCoefficient))
		if interval > e.opt.MaxInterval {
			interval = e.opt.MaxInterval
		}
	}
}

func (e *VideoEncoder) EncodeVideo(ctx context.Context, job *EncodingJob) error {
	tempInputPath := fmt.Sprintf("%s/%s", e.tempPath, job.ObjectName)
	tempOutputDir := fmt.Sprintf("%s/%s_output", e.tempPath, job.RunID)

	// 1. Create output directory
	if err := os.MkdirAll(tempOutputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create temp output dir: %w", err)
	}

	// 2. Download from MinIO
	err := e.downloadFromMinIO(ctx, job.BucketName, job.ObjectName, tempInputPath)
	if err != nil {
		return fmt.Errorf("failed to download from MinIO: %w", err)
	}

	// 3. Encode to HLS
	err = e.EncodeVideoToHLS(ctx, tempInputPath, tempOutputDir)
	if err != nil {
		return fmt.Errorf("HLS encoding failed: %w", err)
	}

	// 4. Upload output directory to MinIO under "video/{runID}/"
	err = e.uploadHLSOutput(ctx, job.BucketName, tempOutputDir, job.VideoID)
	if err != nil {
		return fmt.Errorf("failed to upload HLS output: %w", err)
	}

	// 5. Clean up temp files (optional)
	os.Remove(tempInputPath)
	os.RemoveAll(tempOutputDir)

	return nil
}

func (e *VideoEncoder) EncodeHLS(ctx context.Context, inputPath, outputDir string, res *Resolution) error {
	outputPath := fmt.Sprintf("%s/%s.m3u8", outputDir, res.label)
	segmentPath := fmt.Sprintf("%s/%s_%%03d.ts", outputDir, res.label)

	args := []string{
		"-i", inputPath,
		"-vf", fmt.Sprintf("scale=w=%d:h=%d:force_original_aspect_ratio=decrease", res.width, res.height),
		"-c:a", "aac",
		"-ar", "48000",
		"-c:v", "h264",
		"-profile:v", "main",
		"-crf", "20",
		"-sc_threshold", "0",
		"-g", "48",
		"-keyint_min", "48",
		"-hls_time", "4",
		"-hls_playlist_type", "vod",
		"-b:v", res.bitrate,
		"-maxrate", calculateMaxRate(res.bitrate),
		"-bufsize", calculateBufSize(res.bitrate),
		"-hls_segment_filename", segmentPath,
		outputPath,
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (e *VideoEncoder) EncodeVideoToHLS(ctx context.Context, inputPath, outputDir string) error {
	width, height, err := GetVideoResolution(inputPath)
	if err != nil {
		return fmt.Errorf("ffprobe error: %w", err)
	}

	resolutions := GetResolutionsUnder(width, height)

	var wg sync.WaitGroup
	errCh := make(chan error, len(resolutions))
	sem := make(chan struct{}, e.opt.MaxWorker)

	for _, res := range resolutions {
		res := res
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer func() { <-sem }()
			defer wg.Done()

			ctx := context.Background()

			var attempt int
			interval := e.opt.InitialInterval

			for {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				default:
				}

				fmt.Printf("Encoding %s attempt %d...\n", res.label, attempt+1)
				encodeErr := e.EncodeHLS(ctx, inputPath, outputDir, res)
				if encodeErr == nil {
					errCh <- nil
					return
				}

				attempt++
				if e.opt.MaxAttempts > 0 && attempt >= e.opt.MaxAttempts {
					errCh <- fmt.Errorf("failed at %s after %d attempts: %w", res.label, attempt, encodeErr)
					return
				}

				fmt.Printf("Retrying %s in %v due to error: %v\n", res.label, interval, encodeErr)
				time.Sleep(interval)

				interval = time.Duration(float64(interval) * e.opt.BackoffCoefficient)
				if interval > e.opt.MaxInterval {
					interval = e.opt.MaxInterval
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var combinedErr error
	for err := range errCh {
		if err != nil && combinedErr == nil {
			combinedErr = err
		}
	}
	return combinedErr
}

func (e *VideoEncoder) downloadFromMinIO(ctx context.Context, bucket, object, dest string) error {
	objectReader, err := e.minioClient.GetObject(ctx, bucket, object, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer objectReader.Close()

	localFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, objectReader)
	return err
}

func (e *VideoEncoder) uploadHLSOutput(ctx context.Context, bucket, dirPath, videoID string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", dirPath, file.Name())
		objectName := fmt.Sprintf("video/%s/%s", videoID, file.Name())

		f, err := os.Open(filePath)
		if err != nil {
			return err
		}

		stat, err := f.Stat()
		if err != nil {
			f.Close()
			return err
		}

		_, err = e.minioClient.PutObject(ctx, bucket, objectName, f, stat.Size(), minio.PutObjectOptions{
			ContentType: getContentType(file.Name()),
		})
		f.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func calculateMaxRate(bitrate string) string {
	val, unit := parseBitrate(bitrate)
	return fmt.Sprintf("%d%s", int(float64(val)*1.07), unit)
}

func calculateBufSize(bitrate string) string {
	val, unit := parseBitrate(bitrate)
	return fmt.Sprintf("%d%s", val*2, unit) // typically 2x bitrate
}

func parseBitrate(bitrate string) (int, string) {
	var val int
	var unit string
	fmt.Sscanf(bitrate, "%d%s", &val, &unit)
	return val, unit
}

func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream"
	}
	return mimeType
}
