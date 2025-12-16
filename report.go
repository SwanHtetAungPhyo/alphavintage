package alphavintage

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jung-kurt/gofpdf"
)

var imageCounter int64

// ReportBuilder creates PDF reports with charts and text
type ReportBuilder struct {
	pdf        *gofpdf.Fpdf
	pageWidth  float64
	pageHeight float64
	margin     float64
	fontFamily string
}

// ReportOptions configures the PDF report
type ReportOptions struct {
	Title       string
	Author      string
	Subject     string
	PageSize    string  // A4, Letter, Legal
	Orientation string  // P (portrait), L (landscape)
	MarginMM    float64
}

// DefaultReportOptions returns default report options
func DefaultReportOptions() ReportOptions {
	return ReportOptions{
		Title:       "Stock Analysis Report",
		Author:      "Alpha Vantage Go Client",
		Subject:     "Financial Analysis",
		PageSize:    "A4",
		Orientation: "P",
		MarginMM:    20,
	}
}

// NewReportBuilder creates a new PDF report builder
func NewReportBuilder(opts ReportOptions) *ReportBuilder {
	if opts.PageSize == "" {
		opts.PageSize = "A4"
	}
	if opts.Orientation == "" {
		opts.Orientation = "P"
	}
	if opts.MarginMM == 0 {
		opts.MarginMM = 20
	}

	pdf := gofpdf.New(opts.Orientation, "mm", opts.PageSize, "")
	pdf.SetTitle(opts.Title, true)
	pdf.SetAuthor(opts.Author, true)
	pdf.SetSubject(opts.Subject, true)
	pdf.SetCreationDate(time.Now())
	pdf.SetMargins(opts.MarginMM, opts.MarginMM, opts.MarginMM)
	pdf.SetAutoPageBreak(true, opts.MarginMM+10)

	w, h := pdf.GetPageSize()

	return &ReportBuilder{
		pdf:        pdf,
		pageWidth:  w,
		pageHeight: h,
		margin:     opts.MarginMM,
		fontFamily: "Helvetica",
	}
}

// AddPage adds a new page to the report
func (rb *ReportBuilder) AddPage() *ReportBuilder {
	rb.pdf.AddPage()
	return rb
}

// AddTitle adds a large centered title
func (rb *ReportBuilder) AddTitle(text string) *ReportBuilder {
	rb.pdf.SetFont(rb.fontFamily, "B", 28)
	rb.pdf.SetTextColor(0, 51, 102)
	rb.pdf.SetX(rb.margin)
	rb.pdf.MultiCell(rb.contentWidth(), 14, text, "", "C", false)
	rb.pdf.Ln(8)
	return rb
}

// AddSubtitle adds a subtitle
func (rb *ReportBuilder) AddSubtitle(text string) *ReportBuilder {
	rb.pdf.SetFont(rb.fontFamily, "", 18)
	rb.pdf.SetTextColor(80, 80, 80)
	rb.pdf.SetX(rb.margin)
	rb.pdf.MultiCell(rb.contentWidth(), 9, text, "", "C", false)
	rb.pdf.Ln(5)
	return rb
}

// AddHeading adds a section heading with underline
func (rb *ReportBuilder) AddHeading(text string) *ReportBuilder {
	rb.pdf.Ln(3)
	rb.pdf.SetFont(rb.fontFamily, "B", 16)
	rb.pdf.SetTextColor(0, 82, 147)
	rb.pdf.SetX(rb.margin)
	rb.pdf.CellFormat(rb.contentWidth(), 10, text, "", 1, "L", false, 0, "")
	
	// Add underline
	rb.pdf.SetDrawColor(0, 82, 147)
	rb.pdf.SetLineWidth(0.5)
	y := rb.pdf.GetY()
	rb.pdf.Line(rb.margin, y, rb.margin+rb.contentWidth(), y)
	rb.pdf.Ln(6)
	return rb
}

// AddText adds regular paragraph text
func (rb *ReportBuilder) AddText(text string) *ReportBuilder {
	rb.pdf.SetFont(rb.fontFamily, "", 11)
	rb.pdf.SetTextColor(40, 40, 40)
	rb.pdf.SetX(rb.margin)
	rb.pdf.MultiCell(rb.contentWidth(), 6, text, "", "J", false)
	rb.pdf.Ln(4)
	return rb
}

// AddBoldText adds bold text
func (rb *ReportBuilder) AddBoldText(text string) *ReportBuilder {
	rb.pdf.SetFont(rb.fontFamily, "B", 11)
	rb.pdf.SetTextColor(40, 40, 40)
	rb.pdf.SetX(rb.margin)
	rb.pdf.MultiCell(rb.contentWidth(), 6, text, "", "L", false)
	rb.pdf.Ln(3)
	return rb
}

// AddItalicText adds italic text
func (rb *ReportBuilder) AddItalicText(text string) *ReportBuilder {
	rb.pdf.SetFont(rb.fontFamily, "I", 11)
	rb.pdf.SetTextColor(100, 100, 100)
	rb.pdf.SetX(rb.margin)
	rb.pdf.MultiCell(rb.contentWidth(), 6, text, "", "J", false)
	rb.pdf.Ln(3)
	return rb
}

// AddLineBreak adds vertical space
func (rb *ReportBuilder) AddLineBreak(mm float64) *ReportBuilder {
	rb.pdf.Ln(mm)
	return rb
}

// AddHorizontalLine adds a horizontal separator line
func (rb *ReportBuilder) AddHorizontalLine() *ReportBuilder {
	rb.pdf.SetDrawColor(180, 180, 180)
	rb.pdf.SetLineWidth(0.3)
	y := rb.pdf.GetY() + 2
	rb.pdf.Line(rb.margin, y, rb.pageWidth-rb.margin, y)
	rb.pdf.Ln(6)
	return rb
}

// AddBulletPoint adds a bullet point item
func (rb *ReportBuilder) AddBulletPoint(text string) *ReportBuilder {
	rb.pdf.SetFont(rb.fontFamily, "", 11)
	rb.pdf.SetTextColor(40, 40, 40)
	rb.pdf.SetX(rb.margin)
	rb.pdf.CellFormat(6, 6, "-", "", 0, "L", false, 0, "")
	rb.pdf.MultiCell(rb.contentWidth()-6, 6, text, "", "L", false)
	return rb
}

// AddNumberedItem adds a numbered list item
func (rb *ReportBuilder) AddNumberedItem(num int, text string) *ReportBuilder {
	rb.pdf.SetFont(rb.fontFamily, "", 11)
	rb.pdf.SetTextColor(40, 40, 40)
	rb.pdf.SetX(rb.margin)
	rb.pdf.CellFormat(10, 6, fmt.Sprintf("%d.", num), "", 0, "L", false, 0, "")
	rb.pdf.MultiCell(rb.contentWidth()-10, 6, text, "", "L", false)
	return rb
}

// AddKeyValue adds a key-value pair in a clean format
func (rb *ReportBuilder) AddKeyValue(key, value string) *ReportBuilder {
	rb.pdf.SetX(rb.margin)
	rb.pdf.SetFont(rb.fontFamily, "B", 11)
	rb.pdf.SetTextColor(60, 60, 60)
	rb.pdf.CellFormat(55, 7, key+":", "", 0, "L", false, 0, "")
	rb.pdf.SetFont(rb.fontFamily, "", 11)
	rb.pdf.SetTextColor(40, 40, 40)
	rb.pdf.CellFormat(rb.contentWidth()-55, 7, value, "", 1, "L", false, 0, "")
	return rb
}

// AddTable adds a formatted table
func (rb *ReportBuilder) AddTable(headers []string, rows [][]string) *ReportBuilder {
	if len(headers) == 0 {
		return rb
	}

	colWidth := rb.contentWidth() / float64(len(headers))

	// Header row
	rb.pdf.SetFont(rb.fontFamily, "B", 10)
	rb.pdf.SetFillColor(0, 82, 147)
	rb.pdf.SetTextColor(255, 255, 255)
	rb.pdf.SetX(rb.margin)
	for _, h := range headers {
		rb.pdf.CellFormat(colWidth, 8, h, "1", 0, "C", true, 0, "")
	}
	rb.pdf.Ln(-1)

	// Data rows
	rb.pdf.SetFont(rb.fontFamily, "", 10)
	rb.pdf.SetTextColor(40, 40, 40)
	for i, row := range rows {
		if i%2 == 0 {
			rb.pdf.SetFillColor(245, 245, 245)
		} else {
			rb.pdf.SetFillColor(255, 255, 255)
		}
		rb.pdf.SetX(rb.margin)
		for j, cell := range row {
			if j < len(headers) {
				rb.pdf.CellFormat(colWidth, 7, cell, "1", 0, "C", true, 0, "")
			}
		}
		rb.pdf.Ln(-1)
	}
	rb.pdf.Ln(5)
	return rb
}

func (rb *ReportBuilder) contentWidth() float64 {
	return rb.pageWidth - (2 * rb.margin)
}

func (rb *ReportBuilder) checkPageBreak(height float64) {
	if rb.pdf.GetY()+height > rb.pageHeight-rb.margin-15 {
		rb.pdf.AddPage()
	}
}


func (rb *ReportBuilder) addChartImage(data []byte, name string, widthMM, heightMM float64) {
	// Generate unique name for each image
	uniqueName := fmt.Sprintf("%s_%d", name, atomic.AddInt64(&imageCounter, 1))
	
	reader := bytes.NewReader(data)
	rb.pdf.RegisterImageOptionsReader(uniqueName, gofpdf.ImageOptions{ImageType: "PNG"}, reader)

	// Check if we need a new page
	rb.checkPageBreak(heightMM + 10)

	// Center the image
	x := rb.margin + (rb.contentWidth()-widthMM)/2
	if x < rb.margin {
		x = rb.margin
	}

	rb.pdf.ImageOptions(uniqueName, x, rb.pdf.GetY(), widthMM, heightMM, false,
		gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	rb.pdf.SetY(rb.pdf.GetY() + heightMM + 5)
}

// AddDailyPriceChart generates and adds a price chart
func (rb *ReportBuilder) AddDailyPriceChart(data *TimeSeriesDailyResponse, opts ChartOptions) *ReportBuilder {
	if data == nil || len(data.TimeSeries) == 0 {
		return rb
	}

	if opts.Width == 0 {
		opts.Width = 1000
	}
	if opts.Height == 0 {
		opts.Height = 500
	}

	var buf bytes.Buffer
	if err := GenerateDailyPriceChart(data, &buf, opts); err != nil {
		rb.AddText(fmt.Sprintf("Error generating chart: %v", err))
		return rb
	}

	// Scale to fit page width while maintaining aspect ratio
	imgWidth := rb.contentWidth()
	imgHeight := imgWidth * float64(opts.Height) / float64(opts.Width)
	
	rb.addChartImage(buf.Bytes(), "price", imgWidth, imgHeight)
	return rb
}

// AddCandlestickChart generates and adds a candlestick chart
func (rb *ReportBuilder) AddCandlestickChart(data *TimeSeriesDailyResponse, opts ChartOptions) *ReportBuilder {
	if data == nil || len(data.TimeSeries) == 0 {
		return rb
	}

	if opts.Width == 0 {
		opts.Width = 1000
	}
	if opts.Height == 0 {
		opts.Height = 500
	}

	var buf bytes.Buffer
	if err := GenerateCandlestickChart(data, &buf, opts); err != nil {
		rb.AddText(fmt.Sprintf("Error generating chart: %v", err))
		return rb
	}

	imgWidth := rb.contentWidth()
	imgHeight := imgWidth * float64(opts.Height) / float64(opts.Width)
	
	rb.addChartImage(buf.Bytes(), "candle", imgWidth, imgHeight)
	return rb
}

// AddEarningsChart generates and adds an earnings chart
func (rb *ReportBuilder) AddEarningsChart(data *EarningsResponse, opts ChartOptions) *ReportBuilder {
	if data == nil || len(data.AnnualEarnings) == 0 {
		return rb
	}

	if opts.Width == 0 {
		opts.Width = 800
	}
	if opts.Height == 0 {
		opts.Height = 400
	}

	var buf bytes.Buffer
	if err := GenerateEarningsChart(data, &buf, opts); err != nil {
		rb.AddText(fmt.Sprintf("Error generating chart: %v", err))
		return rb
	}

	imgWidth := rb.contentWidth() * 0.85
	imgHeight := imgWidth * float64(opts.Height) / float64(opts.Width)
	
	rb.addChartImage(buf.Bytes(), "earnings", imgWidth, imgHeight)
	return rb
}

// AddCashFlowChart generates and adds a cash flow chart
func (rb *ReportBuilder) AddCashFlowChart(data *CashFlowResponse, opts ChartOptions) *ReportBuilder {
	if data == nil || len(data.AnnualReports) == 0 {
		return rb
	}

	if opts.Width == 0 {
		opts.Width = 900
	}
	if opts.Height == 0 {
		opts.Height = 450
	}

	var buf bytes.Buffer
	if err := GenerateCashFlowChart(data, &buf, opts); err != nil {
		rb.AddText(fmt.Sprintf("Error generating chart: %v", err))
		return rb
	}

	imgWidth := rb.contentWidth()
	imgHeight := imgWidth * float64(opts.Height) / float64(opts.Width)
	
	rb.addChartImage(buf.Bytes(), "cashflow", imgWidth, imgHeight)
	return rb
}

// AddComparisonChart generates and adds a comparison chart
func (rb *ReportBuilder) AddComparisonChart(datasets map[string]*TimeSeriesDailyResponse, opts ChartOptions) *ReportBuilder {
	if len(datasets) == 0 {
		return rb
	}

	if opts.Width == 0 {
		opts.Width = 1000
	}
	if opts.Height == 0 {
		opts.Height = 500
	}

	var buf bytes.Buffer
	if err := GenerateComparisonChart(datasets, &buf, opts); err != nil {
		rb.AddText(fmt.Sprintf("Error generating chart: %v", err))
		return rb
	}

	imgWidth := rb.contentWidth()
	imgHeight := imgWidth * float64(opts.Height) / float64(opts.Width)
	
	rb.addChartImage(buf.Bytes(), "compare", imgWidth, imgHeight)
	return rb
}

// AddImageFromFile adds an existing PNG image file
func (rb *ReportBuilder) AddImageFromFile(filepath string, widthMM float64) *ReportBuilder {
	if widthMM == 0 {
		widthMM = rb.contentWidth()
	}
	
	x := rb.margin + (rb.contentWidth()-widthMM)/2
	rb.pdf.ImageOptions(filepath, x, rb.pdf.GetY(), widthMM, 0, false,
		gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	rb.pdf.Ln(10)
	return rb
}


// AddMarketStatusSummary adds a formatted market status table
func (rb *ReportBuilder) AddMarketStatusSummary(data *MarketStatusResponse) *ReportBuilder {
	if data == nil || len(data.Markets) == 0 {
		return rb
	}

	var rows [][]string
	for _, m := range data.Markets {
		rows = append(rows, []string{
			m.Region,
			m.MarketType,
			m.CurrentStatus,
			m.LocalOpen + " - " + m.LocalClose,
		})
	}

	rb.AddTable([]string{"Region", "Type", "Status", "Trading Hours"}, rows)
	return rb
}

// AddBalanceSheetSummary adds balance sheet key metrics
func (rb *ReportBuilder) AddBalanceSheetSummary(data *BalanceSheetResponse) *ReportBuilder {
	if data == nil || len(data.AnnualReports) == 0 {
		return rb
	}

	report := data.AnnualReports[0]
	
	rb.AddKeyValue("Fiscal Date", report.FiscalDateEnding)
	rb.AddKeyValue("Total Assets", formatCurrency(report.TotalAssets))
	rb.AddKeyValue("Total Liabilities", formatCurrency(report.TotalLiabilities))
	rb.AddKeyValue("Shareholder Equity", formatCurrency(report.TotalShareholderEquity))
	rb.AddKeyValue("Cash & Equivalents", formatCurrency(report.CashAndCashEquivalentsAtCarryingValue))
	rb.AddKeyValue("Long Term Debt", formatCurrency(report.LongTermDebt))
	rb.AddKeyValue("Goodwill", formatCurrency(report.Goodwill))
	rb.pdf.Ln(5)
	return rb
}

// AddEarningsSummary adds earnings table
func (rb *ReportBuilder) AddEarningsSummary(data *EarningsResponse, years int) *ReportBuilder {
	if data == nil || len(data.AnnualEarnings) == 0 {
		return rb
	}

	if years <= 0 || years > len(data.AnnualEarnings) {
		years = len(data.AnnualEarnings)
	}
	if years > 10 {
		years = 10
	}

	var rows [][]string
	for i := 0; i < years; i++ {
		e := data.AnnualEarnings[i]
		rows = append(rows, []string{e.FiscalDateEnding, "$" + e.ReportedEPS})
	}

	rb.AddTable([]string{"Fiscal Year End", "Earnings Per Share"}, rows)
	return rb
}

// AddCashFlowSummary adds cash flow key metrics
func (rb *ReportBuilder) AddCashFlowSummary(data *CashFlowResponse) *ReportBuilder {
	if data == nil || len(data.AnnualReports) == 0 {
		return rb
	}

	report := data.AnnualReports[0]
	
	rb.AddKeyValue("Fiscal Date", report.FiscalDateEnding)
	rb.AddKeyValue("Operating Cash Flow", formatCurrency(report.OperatingCashflow))
	rb.AddKeyValue("Investing Cash Flow", formatCurrency(report.CashflowFromInvestment))
	rb.AddKeyValue("Financing Cash Flow", formatCurrency(report.CashflowFromFinancing))
	rb.AddKeyValue("Net Income", formatCurrency(report.NetIncome))
	rb.AddKeyValue("Capital Expenditures", formatCurrency(report.CapitalExpenditures))
	rb.AddKeyValue("Dividend Payout", formatCurrency(report.DividendPayout))
	rb.pdf.Ln(5)
	return rb
}

// AddTimestamp adds generation timestamp
func (rb *ReportBuilder) AddTimestamp() *ReportBuilder {
	rb.pdf.SetFont(rb.fontFamily, "I", 10)
	rb.pdf.SetTextColor(120, 120, 120)
	rb.pdf.SetX(rb.margin)
	rb.pdf.CellFormat(rb.contentWidth(), 6,
		fmt.Sprintf("Generated: %s", time.Now().Format("January 2, 2006 at 3:04 PM MST")),
		"", 1, "C", false, 0, "")
	rb.pdf.Ln(3)
	return rb
}

// AddPageNumbers enables page numbering in footer
func (rb *ReportBuilder) AddPageNumbers() *ReportBuilder {
	rb.pdf.SetFooterFunc(func() {
		rb.pdf.SetY(-12)
		rb.pdf.SetFont(rb.fontFamily, "I", 9)
		rb.pdf.SetTextColor(150, 150, 150)
		rb.pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", rb.pdf.PageNo()), "", 0, "C", false, 0, "")
	})
	return rb
}

// AddHeader adds a header to all pages
func (rb *ReportBuilder) AddHeader(text string) *ReportBuilder {
	rb.pdf.SetHeaderFunc(func() {
		rb.pdf.SetY(5)
		rb.pdf.SetFont(rb.fontFamily, "I", 9)
		rb.pdf.SetTextColor(150, 150, 150)
		rb.pdf.CellFormat(0, 10, text, "", 0, "C", false, 0, "")
		rb.pdf.Ln(10)
	})
	return rb
}

// Save writes the PDF to a file
func (rb *ReportBuilder) Save(filename string) error {
	return rb.pdf.OutputFileAndClose(filename)
}

// SaveToBytes returns the PDF as bytes
func (rb *ReportBuilder) SaveToBytes() ([]byte, error) {
	var buf bytes.Buffer
	err := rb.pdf.Output(&buf)
	return buf.Bytes(), err
}

// GetPDF returns the underlying gofpdf for advanced customization
func (rb *ReportBuilder) GetPDF() *gofpdf.Fpdf {
	return rb.pdf
}

// formatCurrency formats large numbers with B/M/K suffixes
func formatCurrency(value string) string {
	if value == "" || value == "None" {
		return "N/A"
	}

	var num float64
	fmt.Sscanf(value, "%f", &num)

	negative := num < 0
	if negative {
		num = -num
	}

	var result string
	switch {
	case num >= 1e12:
		result = fmt.Sprintf("$%.2fT", num/1e12)
	case num >= 1e9:
		result = fmt.Sprintf("$%.2fB", num/1e9)
	case num >= 1e6:
		result = fmt.Sprintf("$%.2fM", num/1e6)
	case num >= 1e3:
		result = fmt.Sprintf("$%.2fK", num/1e3)
	default:
		result = fmt.Sprintf("$%.2f", num)
	}

	if negative {
		result = "-" + result
	}
	return result
}
