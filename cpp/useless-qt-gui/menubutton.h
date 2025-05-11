#ifndef MENUBUTTON_H
#define MENUBUTTON_H
#include <QtWidgets>


class MenuButton : public QPushButton
{
    Q_OBJECT
public:
    explicit MenuButton(QWidget *parent = nullptr, QString text = "",
                    int type = 0, QString iconPath = "");
    void connectMenuButton(QWidget *parent);
    int type = 0;
    int indexOfContentPage = 0;
protected:
protected slots:
signals:
    void emitMenuButtonWasClicked();
    void addButtonWasClicked(MenuButton*);
    void setContentPanelIndex(int);
    void setMenuPanelIndex(int);
private slots:
    void menuButtonWasClicked();
    void setButtonStyleSheet(int state = 0);
    bool event(QEvent *event);
private:
    void setMenuButton(QString text, int type, QString iconPath);
    bool clicked = false;
};
#endif // MENUBUTTON_H
