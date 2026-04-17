import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import Qt.labs.platform 1.1 as Platform
import QtWebEngine

ApplicationWindow {
    id: window
    width: 1400
    height: 800
    visible: true
    title: "QHugo"

    property string currentDir: Platform.StandardPaths.writableLocation(Platform.StandardPaths.DocumentsLocation)
    property int hugoPort: 1313

    Shortcut {
        sequence: "Ctrl+P"
        onActivated: fuzzyFinder.open()
    }
    Shortcut {
        sequence: "Meta+P"
        onActivated: fuzzyFinder.open()
    }

    onCurrentDirChanged: {
        // Retrieve dynamic port and reset web view
        window.hugoPort = FileController.startHugoServer(currentDir)
        webView.url = ""
        // Delay navigation slightly to let Hugo boot up
        previewTimer.start()
    }

    Timer {
        id: previewTimer
        interval: 1000
        onTriggered: {
            console.log("Loading Hugo Preview on Port:", window.hugoPort)
            webView.url = "http://localhost:" + window.hugoPort
        }
    }

    header: ToolBar {
        RowLayout {
            anchors.fill: parent
            ToolButton {
                text: "Open Hugo Repo"
                onClicked: folderDialog.open()
            }
            ToolButton {
                text: "New Post"
                onClicked: newPostDialog.open()
            }
            Item { Layout.fillWidth: true }
        }
    }

    Platform.FolderDialog {
        id: folderDialog
        onAccepted: window.currentDir = folderDialog.folder
    }

    Dialog {
        id: newPostDialog
        title: "Create New Post"
        standardButtons: Dialog.Ok | Dialog.Cancel
        x: Math.round((parent.width - width) / 2)
        y: Math.round((parent.height - height) / 2)
        
        ColumnLayout {
            Label { text: "Post Title:" }
            TextField {
                id: postTitleField
                Layout.fillWidth: true
                focus: true
            }
        }
        onAccepted: {
            var path = FileController.createPost(window.currentDir, postTitleField.text)
            editor.openFile(path)
            postTitleField.text = ""
        }
    }

    SplitView {
        anchors.fill: parent

        Sidebar {
            id: sidebar
            SplitView.preferredWidth: 250
            SplitView.minimumWidth: 150
            SplitView.maximumWidth: 400
            
            currentDirectory: window.currentDir
            
            onFileSelected: function(path) {
                editor.openFile(path)
            }
            onDirectorySelected: function(path) {
                window.currentDir = path
            }
            onGoUpClicked: {
                window.currentDir = FileController.getParentPath(window.currentDir)
            }
        }

        Editor {
            id: editor
            SplitView.fillWidth: true 
            SplitView.preferredWidth: 500
            repoPath: window.currentDir
            onContentSaved: {
                // Hugo handles live-reloading inside the webview via sockets.
            }
        }

        WebEngineView {
            id: webView
            SplitView.preferredWidth: 600
            SplitView.fillWidth: true
        }
    }

    FuzzyFinder {
        id: fuzzyFinder
        rootPath: window.currentDir
        onFileSelected: function(path) {
            editor.openFile(path)
        }
    }
}
