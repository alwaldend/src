#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <menubutton.h>

namespace Ui {
class MainWindow;
}

class MainWindow : public QMainWindow
{
    Q_OBJECT
public:
    explicit MainWindow(QWidget *parent = nullptr);
    ~MainWindow();
signals:
    void menuButtonWasClicked();
    void destroyAddButton();
    void systemButtonWasClicked();
private slots:
    void emitMenuButtonWasClicked();
    void addButtonWasClicked(MenuButton* button);
    void minimizeApp();
    void maximizeApp();
    void closeApp();
    void setContentPanelIndex(int indexOfContentPage = 0);
    void setMenuPanelIndex(int index);
    void fillMenuPanel();
protected:
    void addNewMenuButton(QString text = "",
                       int type = 0, QString iconPath = "");
    void setContentFrame();
    void setMenuFrame();
    void setMainWindow();
    void addContentPage(MenuButton *button);
private:
    Ui::MainWindow *ui;
    QFrame *menuFrame, *contentFrame;
    QStackedWidget *menuPanel, *contentPanel;
    QVBoxLayout *menuPanelPage_Layout;
    MenuButton *addButton,
               *menuArrowTop, *menuArrowBot;
    int amountOfContentPages = 0, amountOfMenuPages = 0;
};

#endif // MAINWINDOW_H
