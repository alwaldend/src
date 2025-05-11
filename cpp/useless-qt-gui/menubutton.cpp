#include "menubutton.h"
#include <QPixmap>

/*
States:
        0: Normal, 1: Hover, 2: Clicked;
Types:
        0: no icon menu button, 1: custom icon menu button,
        2: home menu button, 3: add menu button,
        4: home menu button not clicked, -51: left arrow, -52: right arrow,
        6: add button without arrows
*/
MenuButton::MenuButton(QWidget *parent, QString text, int type, QString iconPath) : QPushButton(parent)
{
    setMenuButton(text, type, iconPath);
    connectMenuButton(parent);
}

void MenuButton::setMenuButton(QString text, int type, QString iconPath)
{
    setMouseTracking(true);
    setAttribute(Qt::WA_Hover);
    this->type=type;
    this->setFlat(true);
    this->setFont(QFont("Century Gothic", 13, QFont::Normal));
    if (type==2) {
        this->setButtonStyleSheet(2);
        this->clicked=true;
    }
    else this->setButtonStyleSheet();
    if (text==QString("")) {
        this->setFixedSize(40,40);
    }
    else {
        this->setText(text);
        this->setFixedSize(184,40);
    }
    if (type) {
        this->setIcon(QPixmap(iconPath));
        if (type==-51 || type==-52) {
            this->setIconSize(QSize(20,20));
            this->setFixedSize(25,25);
        }
        else {
            this->setIconSize(QSize(40,40));
        }
    }
}

void MenuButton::connectMenuButton(QWidget *parent)
{
    if (this->type==3 || this->type==6) {
        connect(this, SIGNAL(addButtonWasClicked(MenuButton*)),
                parent, SLOT(addButtonWasClicked(MenuButton*)));
    }
    else if (this->type==-51 || this->type==-52) {
        connect(this, SIGNAL(setMenuPanelIndex(int)),
                parent, SLOT(setMenuPanelIndex(int)));
    }
    else {
        connect(this, SIGNAL(emitMenuButtonWasClicked()),
        parent, SLOT(emitMenuButtonWasClicked()));
        connect(parent, SIGNAL(menuButtonWasClicked()),
                this, SLOT(menuButtonWasClicked()));

        connect(this, SIGNAL(setContentPanelIndex(int)),
                parent, SLOT(setContentPanelIndex(int)));
    }
}

void MenuButton::menuButtonWasClicked()
{
    if (this->clicked) {
        setButtonStyleSheet();
        this->clicked=false;
    }
}


void MenuButton::setButtonStyleSheet(int state)
{
    switch(state) {
    case 0:
        if (this->type==-51 || this->type==-52)
            this->setStyleSheet("background-color: rgb(60, 60, 60); border-radius: 10px;"
                                "outline: none;"
                                "border: 1px solid rgb(60, 0, 100); color: white;");
        else
            this->setStyleSheet("background-color: rgb(60, 60, 60); border-radius: 20px;"
                                "outline: none;"
                                "border: 1px solid rgb(60, 0, 100); color: white;");
        break;
    case 1:
        if (this->type==-51 || this->type==-52)
            this->setStyleSheet("background-color: rgb(60, 0, 100); border-radius: 10px;"
                                "outline: none;"
                                "border: 1px solid rgb(60, 0, 100); color: white;");
        else
            this->setStyleSheet("background-color: rgb(60, 0, 100); border-radius: 20px;"
                                "outline: none;"
                                "border: 1px solid rgb(60, 0, 100); color: white;");
        break;
    case 2:
        this->setStyleSheet("background-color: rgb(60, 0, 100); border-radius: 20px;"
                            "border: 1px solid rgb(60, 0, 100); color: white;"
                            "outline: none;"
                            "border-left: 4px solid rgb(217, 0, 255);"
                            "border-right: 4px solid rgb(217, 0, 255);");
        break;
    case 3:
        if (this->type==-51 || this->type==-52)
            this->setStyleSheet("background-color: rgb(60, 0, 100); border-radius: 10px;"
                                "outline: none;"
                                "border: 2px solid rgb(217, 0, 255); color: white;");
        else
            this->setStyleSheet("background-color: rgb(60, 0, 100); border-radius: 20px;"
                                "outline: none;"
                                "border: 2px solid rgb(217, 0, 255); color: white;");
        break;
    }
}


bool MenuButton::event(QEvent *event)
{
    switch(event->type())
    {
    case QEvent::HoverEnter:
        if (!this->clicked)
            setButtonStyleSheet(1);
        return true;
    case QEvent::HoverLeave:
        if (!this->clicked)
            setButtonStyleSheet();
        return true;
    case QEvent::MouseButtonPress:
        setButtonStyleSheet(3);
        return true;
    case QEvent::MouseButtonRelease:
        if (this->type==3 || this->type==6) {
            emit addButtonWasClicked(this);
            setButtonStyleSheet();
        }
        else if (this->type==-51 || this->type==-52) {
            setButtonStyleSheet();
            emit setMenuPanelIndex(this->type);
        }
        else {
            emit setContentPanelIndex(this->indexOfContentPage);
            emit emitMenuButtonWasClicked();
            this->clicked=true;
            setButtonStyleSheet(2);
        }
        return true;
    default:
        break;
    }
    return QWidget::event(event);
}
