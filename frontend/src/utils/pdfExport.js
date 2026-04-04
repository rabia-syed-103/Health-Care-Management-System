import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'

/**
 * Universal PDF export function
 * @param {string} title       - Report title shown at top
 * @param {string[]} headers   - Column headers
 * @param {any[][]} rows       - Array of row arrays (already formatted as strings)
 * @param {string} filename    - Downloaded file name (without .pdf)
 */
export function exportPDF(title, headers, rows, filename) {
  const doc = new jsPDF()

  // ── Header ────────────────────────────────────────────────
  doc.setFontSize(20)
  doc.setTextColor(30, 64, 175)   // blue
  doc.text('Hospital Management System', 14, 18)

  doc.setFontSize(13)
  doc.setTextColor(30, 30, 30)
  doc.text(title, 14, 27)

  doc.setFontSize(9)
  doc.setTextColor(120, 120, 120)
  doc.text(`Generated: ${new Date().toLocaleString()}`, 14, 34)

  // ── Table ─────────────────────────────────────────────────
  autoTable(doc, {
    startY: 40,
    head: [headers],
    body: rows,
    headStyles: {
      fillColor: [30, 64, 175],   // blue header
      textColor: 255,
      fontStyle: 'bold',
      fontSize: 10,
    },
    bodyStyles: {
      fontSize: 9,
      textColor: 30,
    },
    alternateRowStyles: {
      fillColor: [240, 245, 255], // light blue alternate rows
    },
    styles: {
      cellPadding: 3,
      lineColor: [200, 210, 230],
      lineWidth: 0.1,
    },
  })

  // ── Footer ────────────────────────────────────────────────
  const pageCount = doc.internal.getNumberOfPages()
  for (let i = 1; i <= pageCount; i++) {
    doc.setPage(i)
    doc.setFontSize(8)
    doc.setTextColor(150)
    doc.text(
      `Page ${i} of ${pageCount}`,
      doc.internal.pageSize.getWidth() / 2,
      doc.internal.pageSize.getHeight() - 8,
      { align: 'center' }
    )
  }

  doc.save(`${filename}.pdf`)
}