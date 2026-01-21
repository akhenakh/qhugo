import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

Popup {
    id: popup
    width: 600
    height: 400
    modal: true
    focus: true
    parent: Overlay.overlay
    x: Math.round((parent.width - width) / 2)
    y: 100
    closePolicy: Popup.CloseOnEscape | Popup.CloseOnPressOutside

    property string rootPath
    signal fileSelected(string path)

    property var allFiles: []

    onOpened: {
        allFiles = FileController.scanDirectory(rootPath)
        filterField.text = ""
        filterField.forceActiveFocus()
        filter()
    }

    // Helper function to handle the actual selection
    function confirmSelection(filePath) {
        if (filePath !== undefined && filePath !== "") {
            popup.fileSelected(filePath)
            popup.close()
        }
    }

    function filter() {
        resultModel.clear()
        var query = filterField.text.toLowerCase()
        var count = 0
        for (var i = 0; i < allFiles.length; i++) {
            var file = allFiles[i]
            if (query === "" || file.toLowerCase().indexOf(query) !== -1) {
                resultModel.append({ "path": file })
                count++
                if (count > 50) break // limit results for performance
            }
        }
        if (resultList.count > 0) resultList.currentIndex = 0
    }

    background: Rectangle {
        color: "#fff"
        border.color: "#ccc"
        radius: 5
        layer.enabled: true 
    }

    ColumnLayout {
        anchors.fill: parent
        anchors.margins: 10

        TextField {
            id: filterField
            Layout.fillWidth: true
            placeholderText: "Search files..."
            onTextChanged: popup.filter()
            
            Keys.onDownPressed: resultList.incrementCurrentIndex()
            Keys.onUpPressed: resultList.decrementCurrentIndex()
            
            // Fixed Keys Handling for Enter/Return
            Keys.onEnterPressed: {
                if (resultList.count > 0) {
                    popup.confirmSelection(resultModel.get(resultList.currentIndex).path)
                }
            }
            Keys.onReturnPressed: {
                if (resultList.count > 0) {
                    popup.confirmSelection(resultModel.get(resultList.currentIndex).path)
                }
            }
        }

        ListView {
            id: resultList
            Layout.fillWidth: true
            Layout.fillHeight: true
            clip: true
            highlight: Rectangle { color: "#eee"; radius: 3 }
            highlightMoveDuration: 0

            model: ListModel { id: resultModel }
            delegate: ItemDelegate {
                // FIX: Use ListView.view instead of parent to avoid null reference errors
                width: ListView.view.width
                text: path.replace(popup.rootPath, "") // Show relative path
                highlighted: ListView.isCurrentItem
                onClicked: {
                    // Pass the model's 'path' role to the helper
                    popup.confirmSelection(path)
                }
            }
        }
    }
}
