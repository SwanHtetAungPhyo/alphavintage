package alphavintage

import (
	"bytes"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/jung-kurt/gofpdf"
)

var imageCounter int64

// LogoPosition defines where the logo appears
type LogoPosition int

const (
	LogoTopLeft LogoPosition = iota
	LogoTopRight
	LogoTopCenter
)

// ReportBuilder creates PDF reports with charts and text
type ReportBuilder struct {
	pdf        *gofpdf.Fpdf
	pageWidth  float64
	pageHeight float64
	margin     float64
	fontFamily string
	logoPath   string
	logoPos    LogoPosition
	logoWidth  float64
}

// ReportOptions configures the PDF report
type ReportOptions struct {
	Title       string
	Author      string
	Subject     string
	PageSize    string  // A4, Letter, Legal
	Orientation string  // P (portrait), L (landscape)
	MarginMM    float64
	LogoPath    string       // Path to logo PNG file
	LogoPosition LogoPosition // Where to place logo
	LogoWidthMM float64      // Logo width in mm (height auto-calculated)
}

// DefaultReportOptions returns default report options
func DefaultReportOptions() ReportOptions {
	return ReportOptions{
		Title:        "Stock Analysis Report",
		Author:       "Alpha Vantage Go Client",
		Subject:      "Financial Analysis",
		PageSize:     "A4",
		Orientation:  "P",
		MarginMM:     20,
		LogoPosition: LogoTopRight,
		LogoWidthMM:  25,
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
	if opts.LogoWidthMM == 0 {
		opts.LogoWidthMM = 25
	}

	pdf := gofpdf.New(opts.Orientation, "mm", opts.PageSize, "")
	pdf.SetTitle(opts.Title, true)
	pdf.SetAuthor(opts.Author, true)
	pdf.SetSubject(opts.Subject, true)
	pdf.SetCreationDate(time.Now())
	pdf.SetMargins(opts.MarginMM, opts.MarginMM, opts.MarginMM)
	pdf.SetAutoPageBreak(true, opts.MarginMM+10)

	w, h := pdf.GetPageSize()

	rb := &ReportBuilder{
		pdf:        pdf,
		pageWidth:  w,
		pageHeight: h,
		margin:     opts.MarginMM,
		fontFamily: "Helvetica",
		logoPath:   opts.LogoPath,
		logoPos:    opts.LogoPosition,
		logoWidth:  opts.LogoWidthMM,
	}

	// Set up logo in header if provided
	if opts.LogoPath != "" {
		rb.setupLogoHeader()
	}

	return rb
}

// SetLogo sets or updates the logo
func (rb *ReportBuilder) SetLogo(path string, position LogoPosition, widthMM float64) *ReportBuilder {
	rb.logoPath = path
	rb.logoPos = position
	if widthMM > 0 {
		rb.logoWidth = widthMM
	}
	rb.setupLogoHeader()
	return rb
}

func (rb *ReportBuilder) setupLogoHeader() {
	if rb.logoPath == "" {
		return
	}

	logoPath := rb.logoPath
	logoWidth := rb.logoWidth
	logoPos := rb.logoPos
	margin := rb.margin
	pageWidth := rb.pageWidth

	rb.pdf.SetHeaderFuncMode(func() {
		var x float64
		switch logoPos {
		case LogoTopLeft:
			x = margin
		case LogoTopRight:
			x = pageWidth - margin - logoWidth
		case LogoTopCenter:
			x = (pageWidth - logoWidth) / 2
		}
		rb.pdf.ImageOptions(logoPath, x, 5, logoWidth, 0, false,
			gofpdf.ImageOptions{ImageType: "", ReadDpi: true}, 0, "")
	}, true)
}

// sanitizeText cleans text for PDF rendering (fixes encoding issues)
func sanitizeText(text string) string {
	result := text
	
	// Replace smart quotes and apostrophes
	result = strings.ReplaceAll(result, "\u2019", "'")  // right single quote
	result = strings.ReplaceAll(result, "\u2018", "'")  // left single quote
	result = strings.ReplaceAll(result, "\u201C", "\"") // left double quote
	result = strings.ReplaceAll(result, "\u201D", "\"") // right double quote
	result = strings.ReplaceAll(result, "\u0027", "'")  // apostrophe
	result = strings.ReplaceAll(result, "\u00B4", "'")  // acute accent
	result = strings.ReplaceAll(result, "\u2032", "'")  // prime
	result = strings.ReplaceAll(result, "\u2033", "\"") // double prime
	
	// Replace dashes and hyphens
	result = strings.ReplaceAll(result, "\u2013", "-")  // en dash
	result = strings.ReplaceAll(result, "\u2014", "-")  // em dash
	result = strings.ReplaceAll(result, "\u2212", "-")  // minus sign
	result = strings.ReplaceAll(result, "\u2010", "-")  // hyphen
	result = strings.ReplaceAll(result, "\u2011", "-")  // non-breaking hyphen
	
	// Replace special characters
	result = strings.ReplaceAll(result, "\u2026", "...") // ellipsis
	result = strings.ReplaceAll(result, "\u00A0", " ")   // non-breaking space
	result = strings.ReplaceAll(result, "\u2022", "-")   // bullet
	result = strings.ReplaceAll(result, "\u00B7", "-")   // middle dot
	result = strings.ReplaceAll(result, "\u2023", "-")   // triangular bullet
	result = strings.ReplaceAll(result, "\u25E6", "-")   // white bullet
	result = strings.ReplaceAll(result, "\u00A9", "(c)") // copyright
	result = strings.ReplaceAll(result, "\u00AE", "(R)") // registered
	result = strings.ReplaceAll(result, "\u2122", "(TM)")// trademark
	result = strings.ReplaceAll(result, "\u00B0", " deg")// degree
	result = strings.ReplaceAll(result, "\u00D7", "x")   // multiplication
	result = strings.ReplaceAll(result, "\u00F7", "/")   // division
	result = strings.ReplaceAll(result, "\u2248", "~")   // approximately
	result = strings.ReplaceAll(result, "\u2260", "!=")  // not equal
	result = strings.ReplaceAll(result, "\u2264", "<=")  // less than or equal
	result = strings.ReplaceAll(result, "\u2265", ">=")  // greater than or equal
	result = strings.ReplaceAll(result, "\u221E", "inf") // infinity
	
	// Replace currency symbols that may cause issues
	result = strings.ReplaceAll(result, "\u20AC", "EUR") // euro
	result = strings.ReplaceAll(result, "\u00A3", "GBP") // pound
	result = strings.ReplaceAll(result, "\u00A5", "JPY") // yen
	
	// Remove markdown formatting
	result = strings.ReplaceAll(result, "**", "")   // bold
	result = strings.ReplaceAll(result, "__", "")   // bold alt
	result = strings.ReplaceAll(result, "~~", "")   // strikethrough
	result = strings.ReplaceAll(result, "```", "")  // code block
	result = strings.ReplaceAll(result, "`", "")    // inline code
	
	// Clean up markdown headers (### Header -> Header)
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		trimmed := strings.TrimLeft(line, " ")
		if strings.HasPrefix(trimmed, "###") {
			lines[i] = strings.TrimPrefix(trimmed, "###")
			lines[i] = strings.TrimSpace(lines[i])
		} else if strings.HasPrefix(trimmed, "##") {
			lines[i] = strings.TrimPrefix(trimmed, "##")
			lines[i] = strings.TrimSpace(lines[i])
		} else if strings.HasPrefix(trimmed, "#") {
			lines[i] = strings.TrimPrefix(trimmed, "#")
			lines[i] = strings.TrimSpace(lines[i])
		}
		// Clean markdown list items
		if strings.HasPrefix(trimmed, "* ") {
			lines[i] = "- " + strings.TrimPrefix(trimmed, "* ")
		}
	}
	result = strings.Join(lines, "\n")
	
	// Clean up markdown links [text](url) -> text
	for {
		start := strings.Index(result, "[")
		if start == -1 {
			break
		}
		mid := strings.Index(result[start:], "](")
		if mid == -1 {
			break
		}
		end := strings.Index(result[start+mid:], ")")
		if end == -1 {
			break
		}
		linkText := result[start+1 : start+mid]
		result = result[:start] + linkText + result[start+mid+end+1:]
	}
	
	// Clean up markdown tables (basic cleanup)
	result = strings.ReplaceAll(result, "|", " ")
	result = strings.ReplaceAll(result, "---", "")
	
	// Remove any remaining non-ASCII that could cause issues
	var cleaned strings.Builder
	for _, r := range result {
		if r < 128 || r == '\n' || r == '\t' {
			cleaned.WriteRune(r)
		} else if r >= 0x00C0 && r <= 0x00FF {
			// Keep extended Latin characters but map common ones
			switch r {
			case 0x00E0, 0x00E1, 0x00E2, 0x00E3, 0x00E4, 0x00E5:
				cleaned.WriteRune('a')
			case 0x00E8, 0x00E9, 0x00EA, 0x00EB:
				cleaned.WriteRune('e')
			case 0x00EC, 0x00ED, 0x00EE, 0x00EF:
				cleaned.WriteRune('i')
			case 0x00F2, 0x00F3, 0x00F4, 0x00F5, 0x00F6:
				cleaned.WriteRune('o')
			case 0x00F9, 0x00FA, 0x00FB, 0x00FC:
				cleaned.WriteRune('u')
			case 0x00F1:
				cleaned.WriteRune('n')
			case 0x00E7:
				cleaned.WriteRune('c')
			default:
				cleaned.WriteRune(' ')
			}
		} else {
			cleaned.WriteRune(' ')
		}
	}
	
	return cleaned.String()
}

// AddPage adds a new page to the report
func (rb *ReportBuilder) AddPage() *ReportBuilder {
	rb.pdf.AddPage()
	// Add space for logo if present
	if rb.logoPath != "" {
		rb.pdf.Ln(15)
	}
	return rb
}

// AddTitle adds a large centered title
func (rb *ReportBuilder) AddTitle(text string) *ReportBuilder {
	rb.pdf.SetFont(rb.fontFamily, "B", 28)
	rb.pdf.SetTextColor(0, 51, 102)
	rb.pdf.SetX(rb.margin)
	rb.pdf.MultiCell(rb.contentWidth(), 14, sanitizeText(text), "", "C", false)
	rb.pdf.Ln(8)
	return rb
}

// AddSubtitle adds a subtitle
func (rb *ReportBuilder) AddSubtitle(text string) *ReportBuilder {
	rb.pdf.SetFont(rb.fontFamily, "", 18)
	rb.pdf.SetTextColor(80, 80, 80)
	rb.pdf.SetX(rb.margin)
	rb.pdf.MultiCell(rb.contentWidth(), 9, sanitizeText(text), "", "C", false)
	rb.pdf.Ln(5)
	return rb
}

// AddHeading adds a section heading with underline
func (rb *ReportBuilder) AddHeading(text string) *ReportBuilder {
	rb.pdf.Ln(3)
	rb.pdf.SetFont(rb.fontFamily, "B", 16)
	rb.pdf.SetTextColor(0, 82, 147)
	rb.pdf.SetX(rb.margin)
	rb.pdf.CellFormat(rb.contentWidth(), 10, sanitizeText(text), "", 1, "L", false, 0, "")

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
	rb.pdf.MultiCell(rb.contentWidth(), 6, sanitizeText(text), "", "L", false)
	rb.pdf.Ln(4)
	return rb
}

// AddBoldText adds bold text
func (rb *ReportBuilder) AddBoldText(text string) *ReportBuilder {
	rb.pdf.SetFont(rb.fontFamily, "B", 11)
	rb.pdf.SetTextColor(40, 40, 40)
	rb.pdf.SetX(rb.margin)
	rb.pdf.MultiCell(rb.contentWidth(), 6, sanitizeText(text), "", "L", false)
	rb.pdf.Ln(3)
	return rb
}

// AddItalicText adds italic text
func (rb *ReportBuilder) AddItalicText(text string) *ReportBuilder {
	rb.pdf.SetFont(rb.fontFamily, "I", 11)
	rb.pdf.SetTextColor(100, 100, 100)
	rb.pdf.SetX(rb.margin)
	rb.pdf.MultiCell(rb.contentWidth(), 6, sanitizeText(text), "", "L", false)
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
	rb.pdf.MultiCell(rb.contentWidth()-6, 6, sanitizeText(text), "", "L", false)
	return rb
}

// AddNumberedItem adds a numbered list item
func (rb *ReportBuilder) AddNumberedItem(num int, text string) *ReportBuilder {
	rb.pdf.SetFont(rb.fontFamily, "", 11)
	rb.pdf.SetTextColor(40, 40, 40)
	rb.pdf.SetX(rb.margin)
	rb.pdf.CellFormat(10, 6, fmt.Sprintf("%d.", num), "", 0, "L", false, 0, "")
	rb.pdf.MultiCell(rb.contentWidth()-10, 6, sanitizeText(text), "", "L", false)
	return rb
}

// AddKeyValue adds a key-value pair in a clean format
func (rb *ReportBuilder) AddKeyValue(key, value string) *ReportBuilder {
	rb.pdf.SetX(rb.margin)
	rb.pdf.SetFont(rb.fontFamily, "B", 11)
	rb.pdf.SetTextColor(60, 60, 60)
	rb.pdf.CellFormat(55, 7, sanitizeText(key)+":", "", 0, "L", false, 0, "")
	rb.pdf.SetFont(rb.fontFamily, "", 11)
	rb.pdf.SetTextColor(40, 40, 40)
	rb.pdf.CellFormat(rb.contentWidth()-55, 7, sanitizeText(value), "", 1, "L", false, 0, "")
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
		rb.pdf.CellFormat(colWidth, 8, sanitizeText(h), "1", 0, "C", true, 0, "")
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
				rb.pdf.CellFormat(colWidth, 7, sanitizeText(cell), "1", 0, "C", true, 0, "")
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
	uniqueName := fmt.Sprintf("%s_%d", name, atomic.AddInt64(&imageCounter, 1))

	reader := bytes.NewReader(data)
	rb.pdf.RegisterImageOptionsReader(uniqueName, gofpdf.ImageOptions{ImageType: "PNG"}, reader)

	rb.checkPageBreak(heightMM + 10)

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

// AddIntradayChart generates and adds an intraday chart
func (rb *ReportBuilder) AddIntradayChart(data *TimeSeriesIntradayResponse, opts ChartOptions) *ReportBuilder {
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
	if err := GenerateIntradayChart(data, &buf, opts); err != nil {
		rb.AddText(fmt.Sprintf("Error generating chart: %v", err))
		return rb
	}

	imgWidth := rb.contentWidth()
	imgHeight := imgWidth * float64(opts.Height) / float64(opts.Width)
	rb.addChartImage(buf.Bytes(), "intraday", imgWidth, imgHeight)
	return rb
}

// AddIntradaySummary adds intraday summary statistics
func (rb *ReportBuilder) AddIntradaySummary(summary *IntradaySummary) *ReportBuilder {
	if summary == nil {
		return rb
	}
	rb.AddKeyValue("Symbol", summary.Symbol)
	rb.AddKeyValue("Date", summary.Date)
	rb.AddKeyValue("Interval", summary.Interval)
	rb.AddKeyValue("Open", fmt.Sprintf("$%.2f", summary.Open))
	rb.AddKeyValue("High", fmt.Sprintf("$%.2f", summary.High))
	rb.AddKeyValue("Low", fmt.Sprintf("$%.2f", summary.Low))
	rb.AddKeyValue("Close", fmt.Sprintf("$%.2f", summary.Close))
	rb.AddKeyValue("Total Volume", formatVolume(float64(summary.TotalVol)))
	rb.AddKeyValue("Data Points", fmt.Sprintf("%d", summary.DataPoints))
	rb.pdf.Ln(5)
	return rb
}

// AddDailyRangeSummary adds date range summary statistics
func (rb *ReportBuilder) AddDailyRangeSummary(summary *DailyRangeSummary) *ReportBuilder {
	if summary == nil {
		return rb
	}
	rb.AddKeyValue("Symbol", summary.Symbol)
	rb.AddKeyValue("Period", fmt.Sprintf("%s to %s", summary.StartDate, summary.EndDate))
	rb.AddKeyValue("Trading Days", fmt.Sprintf("%d", summary.TradingDays))
	rb.AddLineBreak(3)
	rb.AddKeyValue("Period Open", fmt.Sprintf("$%.2f", summary.PeriodOpen))
	rb.AddKeyValue("Period High", fmt.Sprintf("$%.2f (%s)", summary.PeriodHigh, summary.HighDate))
	rb.AddKeyValue("Period Low", fmt.Sprintf("$%.2f (%s)", summary.PeriodLow, summary.LowDate))
	rb.AddKeyValue("Period Close", fmt.Sprintf("$%.2f", summary.PeriodClose))
	rb.AddLineBreak(3)
	sign := ""
	if summary.PriceChange >= 0 {
		sign = "+"
	}
	rb.AddKeyValue("Price Change", fmt.Sprintf("%s$%.2f (%s%.2f%%)", sign, summary.PriceChange, sign, summary.PriceChangePct))
	rb.AddKeyValue("Total Volume", formatVolume(float64(summary.TotalVolume)))
	rb.AddKeyValue("Avg Daily Volume", formatVolume(float64(summary.AvgVolume)))
	rb.pdf.Ln(5)
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

// AddImageFromFile adds an existing image file
func (rb *ReportBuilder) AddImageFromFile(filepath string, widthMM float64) *ReportBuilder {
	if widthMM == 0 {
		widthMM = rb.contentWidth()
	}
	x := rb.margin + (rb.contentWidth()-widthMM)/2
	rb.pdf.ImageOptions(filepath, x, rb.pdf.GetY(), widthMM, 0, false,
		gofpdf.ImageOptions{ImageType: "", ReadDpi: true}, 0, "")
	rb.pdf.Ln(10)
	return rb
}

// AddLogo adds a logo at a specific position on the current page
func (rb *ReportBuilder) AddLogo(filepath string, position LogoPosition, widthMM float64) *ReportBuilder {
	if widthMM == 0 {
		widthMM = 30
	}
	var x float64
	switch position {
	case LogoTopLeft:
		x = rb.margin
	case LogoTopRight:
		x = rb.pageWidth - rb.margin - widthMM
	case LogoTopCenter:
		x = (rb.pageWidth - widthMM) / 2
	}
	rb.pdf.ImageOptions(filepath, x, rb.margin, widthMM, 0, false,
		gofpdf.ImageOptions{ImageType: "", ReadDpi: true}, 0, "")
	return rb
}


// AddMarketStatusSummary adds a formatted market status table
func (rb *ReportBuilder) AddMarketStatusSummary(data *MarketStatusResponse) *ReportBuilder {
	if data == nil || len(data.Markets) == 0 {
		return rb
	}
	var rows [][]string
	for _, m := range data.Markets {
		rows = append(rows, []string{m.Region, m.MarketType, m.CurrentStatus, m.LocalOpen + " - " + m.LocalClose})
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
		fmt.Sprintf("Generated: %s", time.Now().Format("January 2, 2006 at 3:04 PM")),
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

// AddHeader adds a text header to all pages
func (rb *ReportBuilder) AddHeader(text string) *ReportBuilder {
	rb.pdf.SetHeaderFuncMode(func() {
		rb.pdf.SetY(5)
		rb.pdf.SetFont(rb.fontFamily, "I", 9)
		rb.pdf.SetTextColor(150, 150, 150)
		rb.pdf.CellFormat(0, 10, sanitizeText(text), "", 0, "C", false, 0, "")
		rb.pdf.Ln(10)
	}, true)
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

// AddAISummary adds an AI-generated analysis summary section
func (rb *ReportBuilder) AddAISummary(summary *AnalysisSummary) *ReportBuilder {
	if summary == nil {
		return rb
	}
	if summary.Executive != "" {
		rb.AddHeading("Executive Summary")
		rb.AddText(summary.Executive)
		rb.AddLineBreak(5)
	}
	if summary.PriceAnalysis != "" {
		rb.AddHeading("Price Analysis")
		rb.AddText(summary.PriceAnalysis)
		rb.AddLineBreak(5)
	}
	if summary.Fundamentals != "" {
		rb.AddHeading("Fundamental Analysis")
		rb.AddText(summary.Fundamentals)
		rb.AddLineBreak(5)
	}
	if summary.Risks != "" {
		rb.AddHeading("Risk Assessment")
		rb.AddText(summary.Risks)
		rb.AddLineBreak(5)
	}
	if summary.Outlook != "" {
		rb.AddHeading("Outlook")
		rb.AddText(summary.Outlook)
		rb.AddLineBreak(5)
	}
	return rb
}

// AddAIExecutiveSummary adds just the executive summary
func (rb *ReportBuilder) AddAIExecutiveSummary(summary string) *ReportBuilder {
	if summary == "" {
		return rb
	}
	rb.AddHeading("AI Executive Summary")
	rb.AddText(summary)
	rb.AddLineBreak(5)
	return rb
}

// AddAIInsight adds a custom AI insight with a title
func (rb *ReportBuilder) AddAIInsight(title, content string) *ReportBuilder {
	if content == "" {
		return rb
	}
	rb.AddHeading(title)
	rb.AddText(content)
	rb.AddLineBreak(5)
	return rb
}


// Financial Datasets API Report Methods

// AddFDCompanyInfo adds company information section
func (rb *ReportBuilder) AddFDCompanyInfo(company *FDCompanyFacts) *ReportBuilder {
	if company == nil {
		return rb
	}
	rb.AddHeading(company.Name + " (" + company.Ticker + ")")
	rb.AddKeyValue("Industry", company.Industry)
	rb.AddKeyValue("Sector", company.Sector)
	rb.AddKeyValue("Exchange", company.Exchange)
	rb.AddKeyValue("Location", company.Location)
	rb.AddKeyValue("Employees", fmt.Sprintf("%.0f", company.NumberOfEmployees))
	rb.AddKeyValue("Market Cap", formatLargeNumber(company.MarketCap))
	rb.AddKeyValue("Website", company.WebsiteURL)
	rb.pdf.Ln(5)
	return rb
}

// AddFDPriceSnapshot adds real-time price snapshot
func (rb *ReportBuilder) AddFDPriceSnapshot(snapshot *FDPriceSnapshot) *ReportBuilder {
	if snapshot == nil {
		return rb
	}
	rb.AddKeyValue("Current Price", fmt.Sprintf("$%.2f", snapshot.Price))
	rb.AddKeyValue("Day Change", fmt.Sprintf("$%.2f (%.2f%%)", snapshot.DayChange, snapshot.DayChangePercent))
	rb.AddKeyValue("Market Cap", formatLargeNumber(snapshot.MarketCap))
	rb.AddKeyValue("As of", snapshot.Time)
	rb.pdf.Ln(5)
	return rb
}

// AddFDIncomeStatementSummary adds income statement summary
func (rb *ReportBuilder) AddFDIncomeStatementSummary(statements []FDIncomeStatement, count int) *ReportBuilder {
	if len(statements) == 0 {
		return rb
	}
	if count <= 0 || count > len(statements) {
		count = len(statements)
	}
	if count > 5 {
		count = 5
	}

	var rows [][]string
	for i := 0; i < count; i++ {
		s := statements[i]
		rows = append(rows, []string{
			s.ReportPeriod,
			formatLargeNumber(s.Revenue),
			formatLargeNumber(s.NetIncome),
			fmt.Sprintf("$%.2f", s.EarningsPerShare),
		})
	}
	rb.AddTable([]string{"Period", "Revenue", "Net Income", "EPS"}, rows)
	return rb
}

// AddFDBalanceSheetSummary adds balance sheet summary
func (rb *ReportBuilder) AddFDBalanceSheetSummary(sheets []FDBalanceSheet) *ReportBuilder {
	if len(sheets) == 0 {
		return rb
	}
	s := sheets[0]
	rb.AddKeyValue("Report Period", s.ReportPeriod)
	rb.AddKeyValue("Total Assets", formatLargeNumber(s.TotalAssets))
	rb.AddKeyValue("Total Liabilities", formatLargeNumber(s.TotalLiabilities))
	rb.AddKeyValue("Shareholders Equity", formatLargeNumber(s.ShareholdersEquity))
	rb.AddKeyValue("Cash & Equivalents", formatLargeNumber(s.CashAndEquivalents))
	rb.AddKeyValue("Total Debt", formatLargeNumber(s.TotalDebt))
	rb.AddKeyValue("Outstanding Shares", formatLargeNumber(s.OutstandingShares))
	rb.pdf.Ln(5)
	return rb
}

// AddFDCashFlowSummary adds cash flow summary
func (rb *ReportBuilder) AddFDCashFlowSummary(statements []FDCashFlowStatement) *ReportBuilder {
	if len(statements) == 0 {
		return rb
	}
	s := statements[0]
	rb.AddKeyValue("Report Period", s.ReportPeriod)
	rb.AddKeyValue("Operating Cash Flow", formatLargeNumber(s.NetCashFlowFromOperations))
	rb.AddKeyValue("Investing Cash Flow", formatLargeNumber(s.NetCashFlowFromInvesting))
	rb.AddKeyValue("Financing Cash Flow", formatLargeNumber(s.NetCashFlowFromFinancing))
	rb.AddKeyValue("Free Cash Flow", formatLargeNumber(s.FreeCashFlow))
	rb.AddKeyValue("Capital Expenditure", formatLargeNumber(s.CapitalExpenditure))
	rb.pdf.Ln(5)
	return rb
}

// AddFDFinancialMetrics adds financial metrics/ratios
func (rb *ReportBuilder) AddFDFinancialMetrics(metrics *FDFinancialMetrics) *ReportBuilder {
	if metrics == nil {
		return rb
	}

	// Valuation
	rb.AddBoldText("Valuation Metrics")
	rb.AddKeyValue("P/E Ratio", fmt.Sprintf("%.2f", metrics.PriceToEarningsRatio))
	rb.AddKeyValue("P/B Ratio", fmt.Sprintf("%.2f", metrics.PriceToBookRatio))
	rb.AddKeyValue("P/S Ratio", fmt.Sprintf("%.2f", metrics.PriceToSalesRatio))
	rb.AddKeyValue("EV/EBITDA", fmt.Sprintf("%.2f", metrics.EVToEBITDA))
	rb.pdf.Ln(3)

	// Profitability
	rb.AddBoldText("Profitability")
	rb.AddKeyValue("Gross Margin", fmt.Sprintf("%.2f%%", metrics.GrossMargin*100))
	rb.AddKeyValue("Operating Margin", fmt.Sprintf("%.2f%%", metrics.OperatingMargin*100))
	rb.AddKeyValue("Net Margin", fmt.Sprintf("%.2f%%", metrics.NetMargin*100))
	rb.AddKeyValue("ROE", fmt.Sprintf("%.2f%%", metrics.ReturnOnEquity*100))
	rb.AddKeyValue("ROA", fmt.Sprintf("%.2f%%", metrics.ReturnOnAssets*100))
	rb.pdf.Ln(3)

	// Liquidity & Leverage
	rb.AddBoldText("Liquidity & Leverage")
	rb.AddKeyValue("Current Ratio", fmt.Sprintf("%.2f", metrics.CurrentRatio))
	rb.AddKeyValue("Quick Ratio", fmt.Sprintf("%.2f", metrics.QuickRatio))
	rb.AddKeyValue("Debt/Equity", fmt.Sprintf("%.2f", metrics.DebtToEquity))
	rb.AddKeyValue("Debt/Assets", fmt.Sprintf("%.2f", metrics.DebtToAssets))
	rb.pdf.Ln(5)
	return rb
}

// AddFDInsiderTrades adds insider trading table
func (rb *ReportBuilder) AddFDInsiderTrades(trades []FDInsiderTrade, count int) *ReportBuilder {
	if len(trades) == 0 {
		return rb
	}
	if count <= 0 || count > len(trades) {
		count = len(trades)
	}
	if count > 10 {
		count = 10
	}

	var rows [][]string
	for i := 0; i < count; i++ {
		t := trades[i]
		txType := "Buy"
		if t.TransactionShares < 0 {
			txType = "Sell"
		}
		rows = append(rows, []string{
			t.TransactionDate,
			t.Name,
			txType,
			fmt.Sprintf("%.0f", abs(t.TransactionShares)),
			formatLargeNumber(abs(t.TransactionValue)),
		})
	}
	rb.AddTable([]string{"Date", "Insider", "Type", "Shares", "Value"}, rows)
	return rb
}

// AddFDInstitutionalOwnership adds institutional ownership table
func (rb *ReportBuilder) AddFDInstitutionalOwnership(ownership []FDInstitutionalOwnership, count int) *ReportBuilder {
	if len(ownership) == 0 {
		return rb
	}
	if count <= 0 || count > len(ownership) {
		count = len(ownership)
	}
	if count > 10 {
		count = 10
	}

	var rows [][]string
	for i := 0; i < count; i++ {
		o := ownership[i]
		rows = append(rows, []string{
			o.Investor,
			formatLargeNumber(o.Shares),
			formatLargeNumber(o.MarketValue),
		})
	}
	rb.AddTable([]string{"Investor", "Shares", "Market Value"}, rows)
	return rb
}

// AddFDNews adds news summary
func (rb *ReportBuilder) AddFDNews(news []FDNews, count int) *ReportBuilder {
	if len(news) == 0 {
		return rb
	}
	if count <= 0 || count > len(news) {
		count = len(news)
	}
	if count > 5 {
		count = 5
	}

	for i := 0; i < count; i++ {
		n := news[i]
		sentiment := n.Sentiment
		if sentiment == "" {
			sentiment = "neutral"
		}
		rb.AddBoldText(n.Date + " - " + n.Source)
		title := n.Title
		if len(title) > 100 {
			title = title[:100] + "..."
		}
		rb.AddText(title + " [" + sentiment + "]")
	}
	rb.pdf.Ln(3)
	return rb
}

// AddFDPriceChart adds price chart from FD data
func (rb *ReportBuilder) AddFDPriceChart(prices []FDPrice, opts ChartOptions) *ReportBuilder {
	if len(prices) == 0 {
		return rb
	}
	if opts.Width == 0 {
		opts.Width = 1000
	}
	if opts.Height == 0 {
		opts.Height = 500
	}

	var buf bytes.Buffer
	if err := GenerateFDPriceChart(prices, &buf, opts); err != nil {
		rb.AddText(fmt.Sprintf("Error generating chart: %v", err))
		return rb
	}

	imgWidth := rb.contentWidth()
	imgHeight := imgWidth * float64(opts.Height) / float64(opts.Width)
	rb.addChartImage(buf.Bytes(), "fd_price", imgWidth, imgHeight)
	return rb
}

// AddFDRevenueChart adds revenue chart from FD data
func (rb *ReportBuilder) AddFDRevenueChart(statements []FDIncomeStatement, opts ChartOptions) *ReportBuilder {
	if len(statements) == 0 {
		return rb
	}
	if opts.Width == 0 {
		opts.Width = 800
	}
	if opts.Height == 0 {
		opts.Height = 400
	}

	var buf bytes.Buffer
	if err := GenerateFDRevenueChart(statements, &buf, opts); err != nil {
		rb.AddText(fmt.Sprintf("Error generating chart: %v", err))
		return rb
	}

	imgWidth := rb.contentWidth() * 0.85
	imgHeight := imgWidth * float64(opts.Height) / float64(opts.Width)
	rb.addChartImage(buf.Bytes(), "fd_revenue", imgWidth, imgHeight)
	return rb
}

// Helper functions
func formatLargeNumber(n float64) string {
	negative := n < 0
	if negative {
		n = -n
	}
	var result string
	switch {
	case n >= 1e12:
		result = fmt.Sprintf("$%.2fT", n/1e12)
	case n >= 1e9:
		result = fmt.Sprintf("$%.2fB", n/1e9)
	case n >= 1e6:
		result = fmt.Sprintf("$%.2fM", n/1e6)
	case n >= 1e3:
		result = fmt.Sprintf("$%.2fK", n/1e3)
	default:
		result = fmt.Sprintf("$%.2f", n)
	}
	if negative {
		result = "-" + result
	}
	return result
}

func abs(n float64) float64 {
	if n < 0 {
		return -n
	}
	return n
}
