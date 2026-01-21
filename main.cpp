#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include "filecontroller.h"
#include "markdownhighlighter.h" 

int main(int argc, char *argv[])
{
    QGuiApplication app(argc, argv);

    // Register the Highlighter Class
    qmlRegisterType<MarkdownHighlighter>("QtMarkdown", 1, 0, "MarkdownHighlighter");

    QQmlApplicationEngine engine;

    FileController fileController;
    engine.rootContext()->setContextProperty("FileController", &fileController);

    const QUrl url(u"qrc:/QtMarkdown/qml/Main.qml"_qs);
    
    QObject::connect(&engine, &QQmlApplicationEngine::objectCreated,
                     &app, [url](QObject *obj, const QUrl &objUrl) {
        if (!obj && url == objUrl)
            QCoreApplication::exit(-1);
    }, Qt::QueuedConnection);
    engine.load(url);

    return app.exec();
}
