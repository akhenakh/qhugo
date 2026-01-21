import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

Item {
    id: root

    function openFile(path) {
        var cleanPath = path.toString().replace("file://", "")
        
        for (var i = 0; i < tabModel.count; ++i) {
            if (tabModel.get(i).filePath === cleanPath) {
                tabBar.currentIndex = i
                return
            }
        }

        var content = FileController.readFile(cleanPath)
        tabModel.append({ "title": cleanPath.split('/').pop(), "filePath": cleanPath, "fileContent": content })
        tabBar.currentIndex = tabModel.count - 1
    }

    ListModel { id: tabModel }

    ColumnLayout {
        anchors.fill: parent
        spacing: 0

        TabBar {
            id: tabBar
            Layout.fillWidth: true
            Repeater {
                model: tabModel
                TabButton {
                    text: title
                    width: implicitWidth + 20
                    onClicked: tabBar.currentIndex = index
                }
            }
        }

        StackLayout {
            currentIndex: tabBar.currentIndex
            Layout.fillWidth: true
            Layout.fillHeight: true

            Repeater {
                model: tabModel
                
                Item {
                    id: tabItem
                    property string path: filePath
                    // Initialize raw content storage with file content
                    property string memoText: fileContent 
                    property bool previewMode: false

                    ColumnLayout {
                        anchors.fill: parent
                        
                        // Toolbar
                        RowLayout {
                            Layout.fillWidth: true
                            Layout.margins: 5
                            Button {
                                text: tabItem.previewMode ? "Edit" : "Preview"
                                onClicked: {
                                    if (!tabItem.previewMode) {
                                        // Switching TO Preview
                                        // 1. Capture raw text from editor
                                        tabItem.memoText = textArea.text
                                        // 2. Change mode (triggers TextFormat change)
                                        tabItem.previewMode = true
                                        // 3. Force re-assignment of text so Qt interprets raw string as Markdown
                                        //    instead of trying to convert the PlainText document structure.
                                        textArea.text = tabItem.memoText
                                    } else {
                                        // Switching TO Edit
                                        // 1. Change mode back to PlainText
                                        tabItem.previewMode = false
                                        // 2. Restore raw text, overwriting any artifacts from the Markdown renderer
                                        textArea.text = tabItem.memoText
                                    }
                                }
                            }
                            Button {
                                text: "Save"
                                onClicked: {
                                    // If in preview mode, save the captured raw text.
                                    // If in edit mode, save the current editor text.
                                    var content = tabItem.previewMode ? tabItem.memoText : textArea.text
                                    FileController.saveFile(path, content)
                                    
                                    // Update memo if saving in edit mode to keep sync
                                    if (!tabItem.previewMode) tabItem.memoText = content
                                }
                            }
                            Item { Layout.fillWidth: true }
                        }

                        // Editor Area
                        ScrollView {
                            id: scrollView
                            Layout.fillWidth: true
                            Layout.fillHeight: true
                            clip: true

                            Row {
                                width: scrollView.availableWidth 
                                
                                // Line Numbers
                                Column {
                                    id: lineNumbers
                                    width: 40
                                    visible: !tabItem.previewMode
                                    Repeater {
                                        model: textArea.lineCount
                                        Label {
                                            width: 40
                                            height: textArea.cursorRectangle.height
                                            horizontalAlignment: Text.AlignRight
                                            padding: 5
                                            text: index + 1
                                            color: "#888"
                                            font: textArea.font
                                        }
                                    }
                                }

                                // Text / Markdown
                                TextArea {
                                    id: textArea
                                    width: parent.width - (lineNumbers.visible ? lineNumbers.width : 0)
                                    
                                    // Initial binding only. 
                                    // Subsequent updates handled by Button logic to prevent conversion data loss.
                                    text: fileContent 
                                    
                                    textFormat: tabItem.previewMode ? TextEdit.MarkdownText : TextEdit.PlainText
                                    
                                    // Styling Logic
                                    font.family: tabItem.previewMode ? Qt.application.font.family : "Courier New"
                                    font.pixelSize: tabItem.previewMode ? 16 : 14
                                    padding: tabItem.previewMode ? 20 : 0
                                    leftPadding: tabItem.previewMode ? 20 : 5
                                    
                                    wrapMode: TextEdit.Wrap
                                    readOnly: tabItem.previewMode
                                    selectByMouse: true
                                    
                                    background: Rectangle {
                                        color: tabItem.previewMode ? "white" : "transparent"
                                        border.width: 0
                                    }
                                    
                                    color: tabItem.previewMode ? "black" : Qt.application.styleHints.colorScheme === Qt.Dark ? "white" : "black"
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}
