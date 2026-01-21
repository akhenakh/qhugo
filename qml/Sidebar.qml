import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import Qt.labs.folderlistmodel 2.15

Item {
    id: root
    // This property receives data from Main.qml. Do not assign to it locally!
    property string currentDirectory
    
    signal fileSelected(string path)
    signal directorySelected(string path) // New signal for folders
    signal goUpClicked() 

    ColumnLayout {
        anchors.fill: parent
        spacing: 0
        
        // Header Area
        Rectangle {
            Layout.fillWidth: true
            height: 40
            color: "#e0e0e0" 
            
            RowLayout {
                anchors.fill: parent
                anchors.margins: 5
                
                Button {
                    icon.name: "go-up"
                    text: ".."
                    Layout.preferredWidth: 40
                    Layout.fillHeight: true
                    
                    // FIX: Use fully qualified enum
                    display: AbstractButton.IconOnly 
                    
                    onClicked: root.goUpClicked()
                    
                    // FIX: Add delay to prevent binding loops and annoyance
                    ToolTip.visible: hovered
                    ToolTip.delay: 500
                    ToolTip.text: "Up one level"
                }

                Label {
                    text: "Files"
                    font.bold: true
                    Layout.fillWidth: true
                    verticalAlignment: Text.AlignVCenter
                    elide: Text.ElideMiddle
                }
            }
        }

        ListView {
            id: list
            Layout.fillWidth: true
            Layout.fillHeight: true
            clip: true

            model: FolderListModel {
                id: folderModel
                folder: "file://" + root.currentDirectory
                showDirsFirst: true
                showDotAndDotDot: false 
                nameFilters: ["*.md", "*.txt", "*.go", "*.cpp", "*.h", "*.qml"]
            }

            delegate: ItemDelegate {
                width: ListView.view.width
                text: fileName
                icon.name: fileIsDir ? "folder" : "text-x-markdown"
                
                onClicked: {
                    if (fileIsDir) {
                        // FIX: Don't assign locally. Emit signal to restore binding.
                        root.directorySelected(filePath)
                    } else {
                        root.fileSelected(filePath)
                    }
                }
            }
        }
    }
}
