package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	rtype "github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/fatih/color"
)

type IdProcessor struct {
	textractClient    *textract.Client
	rekognitionClient *rekognition.Client
}

func NewIdProcessor() *IdProcessor {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		color.Red("✖ Failed to load AWS config: %v", err)
		os.Exit(1)
	}
	color.Green("✔ AWS configuration loaded successfully")
	return &IdProcessor{
		textractClient:    textract.NewFromConfig(cfg),
		rekognitionClient: rekognition.NewFromConfig(cfg),
	}
}

func (p *IdProcessor) printStep(stepNum int, message string) {
	color.Cyan("\nStep %d: %s", stepNum, message)
}

func (p *IdProcessor) DetectFaces(imageBytes []byte) (*rekognition.DetectFacesOutput, error) {
	color.Yellow("🔍 Detecting faces in image...")
	input := &rekognition.DetectFacesInput{
		Image: &rtype.Image{
			Bytes: imageBytes,
		},
		Attributes: []rtype.Attribute{
			rtype.AttributeDefault,
			rtype.AttributeAll,
		},
	}
	return p.rekognitionClient.DetectFaces(context.Background(), input)
}

func (p *IdProcessor) CompareFaces(imageBytes []byte, targetImage []byte) (*rekognition.CompareFacesOutput, error) {
	color.Yellow("🔄 Comparing faces between images...")
	input := &rekognition.CompareFacesInput{
		SourceImage: &rtype.Image{
			Bytes: imageBytes,
		},
		TargetImage: &rtype.Image{
			Bytes: targetImage,
		},
		SimilarityThreshold: aws.Float32(70.0),
	}
	return p.rekognitionClient.CompareFaces(context.Background(), input)
}

func (p *IdProcessor) ProcessStudentId(imageBytes []byte) error {
	p.printStep(1, "Processing ID document with Textract")
	textractInput := &textract.AnalyzeIDInput{
		DocumentPages: []types.Document{
			{Bytes: imageBytes},
		},
	}

	texResult, err := p.textractClient.AnalyzeID(context.Background(), textractInput)
	if err != nil {
		color.Red("✖ Textract analysis failed: %v", err)
		return err
	}
	color.Green("✔ Text extracted successfully")

	color.Cyan("\n📋 Extracted Document Fields:")
	for _, doc := range texResult.IdentityDocuments {
		for _, field := range doc.IdentityDocumentFields {
			color.White("  %-20s: %s", *field.Type.Text, color.YellowString(*field.ValueDetection.Text))
		}
	}

	p.printStep(2, "Analyzing photo with Rekognition")
	faces, err := p.DetectFaces(imageBytes)
	if err != nil {
		color.Red("✖ Face detection failed: %v", err)
		return err
	}
	color.Green("✔ Face analysis completed")

	color.Cyan("\n👤 Face Detection Results:")
	color.White("  Detected %d face(s)", len(faces.FaceDetails))

	if len(faces.FaceDetails) != 1 {
		color.Red("✖ Expected 1 face, found %d", len(faces.FaceDetails))
		return fmt.Errorf("invalid number of faces detected")
	}

	face := faces.FaceDetails[0]
	confidenceThreshold := float32(90.0)

	color.Cyan("\n🧐 Face Quality Analysis:")
	if face.Confidence == nil || *face.Confidence < confidenceThreshold {
		color.Red("✖ Low confidence: %.2f (threshold: %.2f)", *face.Confidence, confidenceThreshold)
		return fmt.Errorf("low face detection confidence")
	}
	color.Green("✔ Confidence: %.2f (threshold: %.2f)", *face.Confidence, confidenceThreshold)

	if face.Quality == nil {
		color.Red("✖ No quality information available")
		return fmt.Errorf("missing face quality data")
	}

	if face.Quality.Brightness == nil || face.Quality.Sharpness == nil {
		color.Red("✖ Incomplete quality metrics")
		return fmt.Errorf("missing quality metrics")
	}

	color.White("  Brightness: %.2f", *face.Quality.Brightness)
	color.White("  Sharpness: %.2f", *face.Quality.Sharpness)

	if *face.Quality.Brightness < 50 || *face.Quality.Sharpness < 50 {
		color.Red("✖ Poor image quality")
		return fmt.Errorf("poor image quality for face detection")
	}
	color.Green("✔ Image quality meets requirements")

	return nil
}

func main() {
	color.Cyan("\n🚀 Starting Student ID Processor")
	color.Cyan("==============================")

	processor := NewIdProcessor()

	processor.printStep(1, "Loading ID image")
	imageBytes, err := os.ReadFile("id.jpg")
	if err != nil {
		color.Red("✖ Error reading ID image: %v", err)
		os.Exit(1)
	}
	color.Green("✔ ID image loaded successfully")

	err = processor.ProcessStudentId(imageBytes)
	if err != nil {
		color.Red("\n✖ ID Processing Failed: %v", err)
		os.Exit(1)
	}

	processor.printStep(3, "Face Verification")
	selfieBytes, err := os.ReadFile("selfie.jpg")
	if err != nil {
		color.Yellow("⚠ Could not read selfie for comparison: %v", err)
		return
	}

	compareResult, err := processor.CompareFaces(imageBytes, selfieBytes)
	if err != nil {
		color.Red("✖ Face comparison failed: %v", err)
		return
	}

	color.Cyan("\n🔍 Face Comparison Results:")
	if len(compareResult.FaceMatches) > 0 {
		similarity := *compareResult.FaceMatches[0].Similarity
		if similarity > 90 {
			color.Green("✔ Strong match: %.2f%% similarity", similarity)
		} else if similarity > 70 {
			color.Yellow("⚠ Moderate match: %.2f%% similarity", similarity)
		} else {
			color.Red("✖ Weak match: %.2f%% similarity", similarity)
		}
	} else {
		color.Red("✖ No face match found")
	}

	color.Green("\n🎉 Processing completed successfully!")
}
