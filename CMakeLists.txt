cmake_minimum_required(VERSION 3.16)
project(QtMarkdown)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(Qt6 REQUIRED COMPONENTS Core Quick Gui QuickControls2)

# Build Go Backend
add_custom_command(
    OUTPUT ${CMAKE_CURRENT_SOURCE_DIR}/backend/backend.a ${CMAKE_CURRENT_SOURCE_DIR}/backend/backend.h
    COMMAND go build -buildmode=c-archive -o backend.a backend.go
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/backend
    COMMENT "Building Go backend..."
)

add_custom_target(go_backend DEPENDS ${CMAKE_CURRENT_SOURCE_DIR}/backend/backend.a)

qt_standard_project_setup()

qt_add_executable(appQtMarkdown
    main.cpp
    filecontroller.cpp
    filecontroller.h
)

qt_add_qml_module(appQtMarkdown
    URI QtMarkdown
    VERSION 1.0
    QML_FILES
        qml/Main.qml
        qml/Editor.qml
        qml/Sidebar.qml
        qml/FuzzyFinder.qml
)

# Link Go Library and Threads (needed for Go runtime)
target_link_libraries(appQtMarkdown
    PRIVATE
    Qt6::Core Qt6::Quick Qt6::Gui Qt6::QuickControls2
    ${CMAKE_CURRENT_SOURCE_DIR}/backend/backend.a
    Threads::Threads
    ${CMAKE_DL_LIBS} 
)

target_include_directories(appQtMarkdown PRIVATE backend)
add_dependencies(appQtMarkdown go_backend)
