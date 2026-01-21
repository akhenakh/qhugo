#include "filecontroller.h"
#include "backend.h" // Still needed for MCP/Go functions
#include <QDirIterator>
#include <QUrl>
#include <QDebug>

FileController::FileController(QObject *parent) : QObject(parent) {
    InitBackend(); // Initialize Go runtime
}

QStringList FileController::scanDirectory(const QString &path) {
    // Convert URL (file:///...) to Local Path (/...)
    QString localPath = path;
    if (path.startsWith("file://")) {
        localPath = QUrl(path).toLocalFile();
    }

    qDebug() << "--- Starting Directory Scan ---";
    qDebug() << "Input Path:" << path;
    qDebug() << "Local Path:" << localPath;

    QStringList fileList;
    
    QDir dir(localPath);
    if (!dir.exists()) {
        qWarning() << "Directory does not exist:" << localPath;
        return fileList;
    }

    QDirIterator it(localPath, QDir::Files | QDir::NoDotAndDotDot, QDirIterator::Subdirectories);
    
    while (it.hasNext()) {
        QString filePath = it.next();

        // Skip heavy folders
        if (filePath.contains("/.git/") || 
            filePath.contains("/node_modules/") || 
            filePath.contains("/build/") ||
            filePath.contains("/.vscode/")) {
            continue;
        }

        fileList.append(filePath);
        
        // Uncomment if you want to see every file (noisy!)
        // qDebug() << "Found:" << filePath;

        if (fileList.size() >= 20000) {
            qWarning() << "Scan limit reached (20000 files)";
            break;
        }
    }
    
    qDebug() << "--- Scan Complete ---";
    qDebug() << "Total files returned:" << fileList.size();

    return fileList;
}


QString FileController::readFile(const QString &path) {
    QString localPath = path;
    if (path.startsWith("file://")) {
        localPath = QUrl(path).toLocalFile();
    }

    char* cPath = localPath.toUtf8().data();
    char* content = ReadFileContent(cPath);
    QString result = QString::fromUtf8(content);
    FreeString(content);
    return result;
}

bool FileController::saveFile(const QString &path, const QString &content) {
    QString localPath = path;
    if (path.startsWith("file://")) {
        localPath = QUrl(path).toLocalFile();
    }
    return SaveFileContent(localPath.toUtf8().data(), content.toUtf8().data()) == 1;
}

QString FileController::connectMcp(const QString &command) {
    char* res = ConnectMCP(command.toUtf8().data());
    QString qRes = QString::fromUtf8(res);
    FreeString(res);
    return qRes;
}

QString FileController::callMcpTool(const QString &name, const QString &argsJson) {
    char* res = CallMCPTool(name.toUtf8().data(), argsJson.toUtf8().data());
    QString qRes = QString::fromUtf8(res);
    FreeString(res);
    return qRes;
}
