import QtQuick
import QtQuick.Controls
import Qt.labs.folderlistmodel 2.15

Item {
    id: root
    property string currentDirectory
    signal fileSelected(string path)

    Column {
        anchors.fill: parent
        
        Label {
            text: "Files"
            font.bold: true
            padding: 10
            width: parent.width
            background: Rectangle { color: "#ddd" }
        }

        ListView {
            id: list
            width: parent.width
            height: parent.height - 30
            clip: true

            model: FolderListModel {
                id: folderModel
                folder: "file://" + root.currentDirectory
                showDirsFirst: true
                nameFilters: ["*.md", "*.txt", "*.go", "*.cpp", "*.h", "*.qml"]
            }

            delegate: ItemDelegate {
                width: parent.width
                text: fileName
                icon.name: fileIsDir ? "folder" : "text-x-markdown"
                
                onClicked: {
                    if (fileIsDir) {
                        // Simple navigation down, requires '..' logic for real app
                        root.currentDirectory = filePath
                    } else {
                        root.fileSelected(filePath)
                    }
                }
            }
        }
    }
}
