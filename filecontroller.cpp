#include "filecontroller.h"
#include "backend.h" 
#include <QDirIterator>
#include <QUrl>
#include <QDebug>
#include <QDateTime>
#include <QRegularExpression>

FileController::FileController(QObject *parent) : QObject(parent) {
    InitBackend(); // Initialize Go runtime
}

FileController::~FileController() {
    StopHugo();
}

QStringList FileController::scanDirectory(const QString &path) {
    QString localPath = path;
    if (path.startsWith("file://")) {
        localPath = QUrl(path).toLocalFile();
    }

    QStringList fileList;
    QDir dir(localPath);
    if (!dir.exists()) {
        return fileList;
    }

    QDirIterator it(localPath, QDir::Files | QDir::NoDotAndDotDot, QDirIterator::Subdirectories);
    
    while (it.hasNext()) {
        QString filePath = it.next();

        // Skip heavy folders and build artifacts
        if (filePath.contains("/.git/") || 
            filePath.contains("/node_modules/") || 
            filePath.contains("/public/") ||
            filePath.contains("/resources/")) {
            continue;
        }

        // Only index markdown files in FuzzyFinder
        if (!filePath.endsWith(".md")) {
            continue;
        }

        fileList.append(filePath);
        
        if (fileList.size() >= 20000) {
            break;
        }
    }
    
    return fileList;
}

QString FileController::getParentPath(const QString &path) {
    QString localPath = path;
    if (path.startsWith("file://")) {
        localPath = QUrl(path).toLocalFile();
    }

    QDir dir(localPath);
    if (dir.cdUp()) {
        return QUrl::fromLocalFile(dir.absolutePath()).toString();
    }
    return path;
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

void FileController::startHugoServer(const QString &repoPath) {
    QString localPath = repoPath;
    if (repoPath.startsWith("file://")) {
        localPath = QUrl(repoPath).toLocalFile();
    }
    StartHugo(localPath.toUtf8().data());
}

void FileController::stopHugoServer() {
    StopHugo();
}

QString FileController::processImage(const QString &srcPath, const QString &repoPath, const QString &docPath) {
    QString localSrc = srcPath;
    if (srcPath.startsWith("file://")) localSrc = QUrl(srcPath).toLocalFile();
    
    QString localRepo = repoPath;
    if (repoPath.startsWith("file://")) localRepo = QUrl(repoPath).toLocalFile();
    
    QString localDoc = docPath;
    if (docPath.startsWith("file://")) localDoc = QUrl(docPath).toLocalFile();

    char* res = ProcessImage(localSrc.toUtf8().data(), localRepo.toUtf8().data(), localDoc.toUtf8().data());
    QString qRes = QString::fromUtf8(res);
    FreeString(res);
    return qRes;
}

QString FileController::createPost(const QString &repoPath, const QString &title) {
    QString localRepo = repoPath;
    if (repoPath.startsWith("file://")) {
        localRepo = QUrl(repoPath).toLocalFile();
    }
    
    QString slug = title.toLower().replace(QRegularExpression("[^a-z0-9]+"), "-");
    QString year = QDateTime::currentDateTime().toString("yyyy");
    
    QString contentDir = localRepo + "/content/post/" + year;
    QString contentPath = contentDir + "/" + slug + ".md";
    
    QString frontmatter = "---\n";
    frontmatter += "title: \"" + title + "\"\n";
    frontmatter += "date: " + QDateTime::currentDateTime().toString(Qt::ISODate) + "\n";
    frontmatter += "draft: true\n";
    frontmatter += "---\n\n";
    
    QDir().mkpath(contentDir);
    saveFile(contentPath, frontmatter);
    return contentPath;
}
