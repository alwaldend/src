#include "systembutton.h"
#include <QPixmap>
#include <QtDebug>
#include <QTimer>
//Types:
//      0: Minimize, 1: Maximize, 2: Close;
//States:
//        0: Normal, 1: Hover, 2: Clicked;
SystemButton::SystemButton(QWidget *parent, int type) : MenuButton(parent)
{
    setSystemButton(type);
    connectSystemButton(parent);
}

void SystemButton::setSystemButton(int type)
{
    setMouseTracking(true);
    setAttribute(Qt::WA_Hover);
    this->type=type;
    this->setFlat(true);
    this->setButtonStyleSheet();
    this->setFixedSize(30,30);
    this->setIconSize(QSize(20,20));
    switchIcon(1);
}

void SystemButton::connectSystemButton(QWidget *parent)
{
    switch(type) {
    case 0:
          connect(this, SIGNAL(systemButtonWasClicked()),
                  parent, SLOT(minimizeApp()));
          connect(parent, SIGNAL(systemButtonWasClicked()),
                  this, SLOT(changeState()));
          break;
    case 1:
          connect(this, SIGNAL(systemButtonWasClicked()),
                  parent, SLOT(maximizeApp()));
          connect(parent, SIGNAL(systemButtonWasClicked()),
                  this, SLOT(changeState()));
          break;
    case 2:
          connect(this, SIGNAL(systemButtonWasClicked()),
                  parent, SLOT(closeApp()));
          connect(parent, SIGNAL(systemButtonWasClicked()),
                  this, SLOT(changeState()));
          break;
    }
}

void SystemButton::setButtonStyleSheet(int state)
{
    switch(state) {
    case 0:
        this->setStyleSheet("border: none; border-radius: 15px;");
        break;
    case 1:
        this->setStyleSheet("border: none; background-color: rgb(60, 0, 100);"
                            "border-radius: 15px;");
        break;
    case 2:
        this->setStyleSheet("background-color: rgb(60, 0, 100); border-radius: 15px;"
                            "border-left: 2px solid rgb(217, 0, 255);"
                            "border-right: 2px solid rgb(217, 0, 255);"
                            "border-top: none; border-bottom: none");
        break;
    case 3:
        this->setStyleSheet("background-color: rgb(60, 0, 100); border-radius: 15px;"
                            "border: 2px solid rgb(217, 0, 255);");
    }
}

void SystemButton::switchIcon(int state)
{
    switch (state) {
    case 0:
            switch (this->type) {
            case 0: this->setIcon(QPixmap(":/img/img/minimize_windowWhite.png")); break;
            case 1: this->setIcon(QPixmap(":/img/img/maximize_windowWhite.png")); break;
            case 2: this->setIcon(QPixmap(":/img/img/close_windiowWhite.png")); break;
            }
            break;
    case 1:
            switch (this->type) {
            case 0: this->setIcon(QPixmap(":/img/img/minimize_window.png")); break;
            case 1: this->setIcon(QPixmap(":/img/img/maximize_window.png")); break;
            case 2: this->setIcon(QPixmap(":/img/img/close_windiow.png")); break;
            }
            break;
    }

}

void SystemButton::changeState()
{
    setButtonStyleSheet();
    switchIcon(1);
    this->clicked=false;
}

bool SystemButton::event(QEvent *event)
{
    switch(event->type())
    {
    case QEvent::HoverEnter:
        if (!this->clicked) {
            setButtonStyleSheet(1);
            switchIcon(0);
        }
        return true;
    case QEvent::HoverLeave:
        if (!this->clicked) {
            switchIcon(1);
            setButtonStyleSheet();
        }
        return true;
    case QEvent::MouseButtonPress:
        setButtonStyleSheet(3);
        switchIcon(0);
        return true;
    case QEvent::MouseButtonRelease:
        this->clicked=true;
        setButtonStyleSheet(2);
        emit systemButtonWasClicked();
        return true;
    default:
        break;
    }
    return QWidget::event(event);
}
