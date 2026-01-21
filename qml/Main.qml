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

    // Use Documents location to avoid scanning the entire Home directory
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
            currentDirectory: window.currentDir
            
            onFileSelected: function(path) {
                editor.openFile(path)
            }
        }

        Editor {
            id: editor
            SplitView.fillWidth: true // Take remaining space
        }
    }

    FuzzyFinder {
        id: fuzzyFinder
        rootPath: window.currentDir
        
        // FIX: Explicitly declare function parameter 'path'
        onFileSelected: function(path) {
            editor.openFile(path)
        }
    }
}
