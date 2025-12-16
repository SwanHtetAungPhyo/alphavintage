package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/SwanHtetAungPhyo/alphavintage"
)

// ByteStreamExample shows how to generate charts as byte streams
// and create PDFs without saving any temporary files
// Run with: go run example/new.go (after commenting out main.go's main)

func ByteStreamExample() {

		client := alphavintage.NewClient("UPS6QRH073V81U5Z")
	

	symbol := "IBM"
	fmt.Printf("Fetching %s data...\n", symbol)

	// Get data
	daily, err := client.GetTimeSeriesDaily(symbol, alphavintage.OutputSizeCompact)
	if err != nil {
		log.Fatal(err)
	}

	// ============================================
	// Method 1: Generate chart as bytes (no file)
	// ============================================
	var chartBuffer bytes.Buffer
	chartOpts := alphavintage.ChartOptions{
		Title:      symbol + " Price Chart",
		Width:      1000,
		Height:     500,
		ShowVolume: true,
	}

	err = alphavintage.GenerateDailyPriceChart(daily, &chartBuffer, chartOpts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Chart generated: %d bytes in memory\n", chartBuffer.Len())

	// You can use chartBuffer.Bytes() for:
	// - Upload to S3
	// - Send via HTTP
	// - Embed in PDF
	// - Store in database

	// ============================================
	// Method 2: Generate PDF with embedded charts (no temp files)
	// ============================================
	report := alphavintage.NewReportBuilder(alphavintage.DefaultReportOptions())
	report.AddPageNumbers()

	report.AddPage()
	report.AddTitle(symbol + " Analysis")
	report.AddTimestamp()

	report.AddPage()
	report.AddHeading("Price Chart")
	report.AddText("Chart generated directly to PDF without temp files.")

	// This internally uses bytes.Buffer - no temp files created
	report.AddDailyPriceChart(daily, chartOpts)

	// ============================================
	// Method 3: Get PDF as bytes (no file saved)
	// ============================================
	pdfBytes, err := report.SaveToBytes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("PDF generated: %d bytes in memory\n", len(pdfBytes))

	// Use pdfBytes for:
	// - Upload to S3 directly
	// - Send as HTTP response
	// - Store in database
	// - Email attachment

	// Example: Upload to S3 (pseudo-code)
	// s3Client.PutObject(ctx, &s3.PutObjectInput{
	//     Bucket: aws.String("my-bucket"),
	//     Key:    aws.String("reports/IBM_report.pdf"),
	//     Body:   bytes.NewReader(pdfBytes),
	//     ContentType: aws.String("application/pdf"),
	// })

	// Or save to file if needed
	err = os.WriteFile(symbol+"_memory_report.pdf", pdfBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Saved to %s_memory_report.pdf\n", symbol)
}
