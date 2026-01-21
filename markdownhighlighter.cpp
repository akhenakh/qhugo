#include "markdownhighlighter.h"
#include <QTextDocument>

MarkdownHighlighter::MarkdownHighlighter(QObject *parent)
    : QSyntaxHighlighter(parent), m_quickDocument(nullptr)
{
    // 1. Headers (Blue) - # Title
    headerFormat.setForeground(Qt::blue);
    headerFormat.setFontWeight(QFont::Bold);
    HighlightingRule rule;
    rule.pattern = QRegularExpression("^(#{1,6})\\s.*");
    rule.format = headerFormat;
    highlightingRules.append(rule);

    // 2. Lists (Red) - * Item or - Item
    listFormat.setForeground(Qt::red);
    rule.pattern = QRegularExpression("^\\s*([*|-])\\s");
    rule.format = listFormat;
    highlightingRules.append(rule);

    // 3. Inline Code (Dark Orange) - `code`
    // Note: DarkOrange is #FF8C00
    codeFormat.setForeground(QColor(0xFF, 0x8C, 0x00)); 
    
    // FIX 1: Use setFontFamilies instead of deprecated setFontFamily
    codeFormat.setFontFamilies(QStringList("Courier New"));
    
    rule.pattern = QRegularExpression("`[^`]+`");
    rule.format = codeFormat;
    highlightingRules.append(rule);
}

void MarkdownHighlighter::highlightBlock(const QString &text)
{
    for (const HighlightingRule &rule : highlightingRules) {
        QRegularExpressionMatchIterator i = rule.pattern.globalMatch(text);
        while (i.hasNext()) {
            QRegularExpressionMatch match = i.next();
            setFormat(match.capturedStart(), match.capturedLength(), rule.format);
        }
    }
}

QQuickTextDocument *MarkdownHighlighter::document() const
{
    return m_quickDocument;
}

void MarkdownHighlighter::setDocument(QQuickTextDocument *document)
{
    if (m_quickDocument == document)
        return;

    m_quickDocument = document;
    
    if (m_quickDocument) {
        // FIX 2: Explicitly call base class method to avoid recursive call to self
        QSyntaxHighlighter::setDocument(m_quickDocument->textDocument());
    } else {
        QSyntaxHighlighter::setDocument(nullptr);
    }

    emit documentChanged();
}
