#ifndef SYSTEMBUTTON_H
#define SYSTEMBUTTON_H
#include <menubutton.h>
#include <QtWidgets>

class SystemButton : public MenuButton
{
    Q_OBJECT
public:
    explicit SystemButton(QWidget *parent = nullptr, int type = 0);
    bool event(QEvent *event);
protected:
    void setButtonStyleSheet(int state = 0);
signals:
    void systemButtonWasClicked();
private slots:
    void changeState();
private:
    void switchIcon(int state = 0);
    void connectSystemButton(QWidget *parent = nullptr);
    void setSystemButton(int type);
    bool clicked = false;
};

#endif // SYSTEMBUTTON_H
