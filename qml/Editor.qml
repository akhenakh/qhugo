import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import QtMarkdown 1.0

Item {
    id: root

    // FIX: Use 'sequences' (list) instead of 'sequence' for StandardKey
    Shortcut {
        sequences: [StandardKey.Close]
        onActivated: closeTab(tabBar.currentIndex)
    }

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

    function closeTab(index) {
        if (index < 0 || index >= tabModel.count) return;
        
        // Remove from model
        tabModel.remove(index);
        
        if (tabBar.currentIndex >= tabModel.count) {
            tabBar.currentIndex = tabModel.count - 1;
        }
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
                    id: tabBtn
                    width: implicitWidth + 30 
                    
                    contentItem: RowLayout {
                        spacing: 5
                        Label {
                            text: title
                            font: tabBtn.font
                            color: tabBtn.checked ? (Qt.application.styleHints.colorScheme === Qt.Dark ? "white" : "black") : "#888"
                            horizontalAlignment: Text.AlignHCenter
                            verticalAlignment: Text.AlignVCenter
                            elide: Text.ElideRight
                            Layout.fillWidth: true
                        }
                        
                        ToolButton {
                            text: "×"
                            font.pixelSize: 16
                            Layout.preferredWidth: 20
                            Layout.preferredHeight: 20
                            background: Rectangle { color: "transparent" }
                            visible: tabBtn.hovered || tabBtn.checked
                            
                            onClicked: {
                                closeTab(index)
                            }
                        }
                    }
                    
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
                                        tabItem.memoText = textArea.text
                                        tabItem.previewMode = true
                                        textArea.text = tabItem.memoText
                                    } else {
                                        tabItem.previewMode = false
                                        textArea.text = tabItem.memoText
                                    }
                                }
                            }
                            Button {
                                text: "Save"
                                onClicked: {
                                    var content = tabItem.previewMode ? tabItem.memoText : textArea.text
                                    FileController.saveFile(path, content)
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

                                TextArea {
                                    id: textArea
                                    width: parent.width - (lineNumbers.visible ? lineNumbers.width : 0)
                                    text: fileContent 
                                    textFormat: tabItem.previewMode ? TextEdit.MarkdownText : TextEdit.PlainText
                                    
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

                                    MarkdownHighlighter {
                                        id: highlighter
                                        document: textArea.textDocument
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}
