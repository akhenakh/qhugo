import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import Qt.labs.folderlistmodel 2.15

Item {
    id: root
    property string currentDirectory
    
    signal fileSelected(string path)
    signal directorySelected(string path)
    signal goUpClicked() 

    ColumnLayout {
        anchors.fill: parent
        spacing: 0
        
        Rectangle {
            Layout.fillWidth: true
            height: 40
            color: "#e0e0e0" 
            
            RowLayout {
                anchors.fill: parent
                anchors.margins: 5
                
                Button {
                    id: goUpBtn
                    icon.name: "go-up"
                    text: ".."
                    Layout.preferredWidth: 40
                    Layout.fillHeight: true
                    display: AbstractButton.IconOnly 
                    onClicked: root.goUpClicked()
                    
                    // FIX: Use explicit ToolTip object to avoid Breeze binding loops
                    ToolTip {
                        parent: goUpBtn
                        visible: goUpBtn.hovered
                        text: "Up one level"
                    }
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
                nameFilters: ["*.md", "*.txt", "*.yaml", "*.toml", "*.json"]
            }

            delegate: ItemDelegate {
                width: ListView.view.width
                text: fileName
                icon.name: fileIsDir ? "folder" : "text-x-markdown"
                
                onClicked: {
                    if (fileIsDir) {
                        root.directorySelected(filePath)
                    } else {
                        root.fileSelected(filePath)
                    }
                }
            }
        }
    }
}
