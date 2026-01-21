import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import Qt.labs.platform 1.1 as Platform

ApplicationWindow {
    id: window
    width: 1200
    height: 800
    visible: true
    title: "QtGo Markdown"

    // Single source of truth for directory state
    property string currentDir: Platform.StandardPaths.writableLocation(Platform.StandardPaths.DocumentsLocation)

    Shortcut {
        sequence: "Ctrl+P"
        onActivated: fuzzyFinder.open()
    }
    Shortcut {
        sequence: "Meta+P"
        onActivated: fuzzyFinder.open()
    }

    SplitView {
        anchors.fill: parent

        Sidebar {
            id: sidebar
            SplitView.preferredWidth: 250
            SplitView.minimumWidth: 150
            SplitView.maximumWidth: 400
            
            // One-way binding: Sidebar follows Window
            currentDirectory: window.currentDir
            
            // Handlers update Window, which flows back to Sidebar
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
