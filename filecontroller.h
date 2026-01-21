#ifndef FILECONTROLLER_H
#define FILECONTROLLER_H

#include <QObject>
#include <QStringList>

class FileController : public QObject
{
    Q_OBJECT
public:
    explicit FileController(QObject *parent = nullptr);

    Q_INVOKABLE QStringList scanDirectory(const QString &path);
    Q_INVOKABLE QString readFile(const QString &path);
    Q_INVOKABLE bool saveFile(const QString &path, const QString &content);
    Q_INVOKABLE QString getParentPath(const QString &path); // ADD THIS
    
    // Go/MCP stuff remains here
    Q_INVOKABLE QString connectMcp(const QString &command);
    Q_INVOKABLE QString callMcpTool(const QString &name, const QString &argsJson);
};

#endif // FILECONTROLLER_H
