#include "mainwindow.h"
#include "ui_mainwindow.h"
#include <QDebug>
#include <QGuiApplication>
#include <QString>
#include <menubutton.h>
#include <systembutton.h>

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent), ui(new Ui::MainWindow) {
  ui->setupUi(this);
  setMainWindow();
  setMenuFrame();
  setContentFrame();
  fillMenuPanel();
}

MainWindow::~MainWindow() { delete ui; }

void MainWindow::setMainWindow() {
  QWidget::setWindowFlags(Qt::FramelessWindowHint);
  ui->Main_Window->setFixedSize(1024, 600);
  QHBoxLayout *mainWindow_Layout = new QHBoxLayout;
  ui->Main_Window->setLayout(mainWindow_Layout);
  mainWindow_Layout->setContentsMargins(0, 0, 0, 0);
  mainWindow_Layout->setSpacing(0);
  ui->Main_Window->setMaximumSize(
      QGuiApplication::primaryScreen()->availableSize());
  QFrame *menuFrame = new QFrame, *contentFrame = new QFrame;
  this->menuFrame = menuFrame;
  this->contentFrame = contentFrame;
  menuFrame->setFixedWidth(224);
  menuFrame->setStyleSheet("background-color: rgb(60, 60, 60);");
  contentFrame->setStyleSheet("border-top:1px solid rgb(217, 0, 255);"
                              "border-left:1px solid rgb(217, 0, 255);");
  contentFrame->setFrameShape(QFrame::Box);
  contentFrame->setFrameShadow(QFrame::Plain);
  mainWindow_Layout->addWidget(menuFrame);
  mainWindow_Layout->addWidget(contentFrame);
}

void MainWindow::setMenuFrame() {
  QVBoxLayout *frame_Layout = new QVBoxLayout;
  this->menuFrame->setLayout(frame_Layout);
  frame_Layout->setContentsMargins(10, 10, 10, 10);
  QStackedWidget *menuPanel = new QStackedWidget;
  this->menuPanel = menuPanel;
  frame_Layout->addWidget(menuPanel);
  menuPanel->setStyleSheet("background-color: rgb(50, 50, 50); "
                           "border-radius: 50px");
  menuPanel->setMaximumHeight(
      QGuiApplication::primaryScreen()->availableSize().height() * 2 / 3);
  QWidget *menuPanelPage = new QWidget;
  menuPanel->addWidget(menuPanelPage);
  QVBoxLayout *menuPanelPage_Layout = new QVBoxLayout;
  menuPanelPage_Layout->setAlignment(Qt::AlignTop);
  menuPanelPage_Layout->setSpacing(35);
  menuPanelPage->setLayout(menuPanelPage_Layout);
  this->menuPanelPage_Layout = menuPanelPage_Layout;
}

void MainWindow::setContentFrame() {
  QVBoxLayout *frameLayout = new QVBoxLayout;
  this->contentFrame->setLayout(frameLayout);
  QHBoxLayout *systemButtonsLayout = new QHBoxLayout;
  frameLayout->addLayout(systemButtonsLayout);
  systemButtonsLayout->setAlignment(Qt::AlignRight);
  systemButtonsLayout->setSpacing(0);
  SystemButton *minimize = new SystemButton(this, 0),
               *maximize = new SystemButton(this, 1),
               *close = new SystemButton(this, 2);
  systemButtonsLayout->addWidget(minimize);
  systemButtonsLayout->addWidget(maximize);
  systemButtonsLayout->addWidget(close);
  QStackedWidget *contentPanel = new QStackedWidget;
  this->contentPanel = contentPanel;
  contentPanel->setStyleSheet("border: none;");
  frameLayout->addWidget(contentPanel);
}

void MainWindow::fillMenuPanel() {
  this->amountOfMenuPages++;
  if (this->menuPanel->currentIndex()) {
    addNewMenuButton("", 4, ":/img/img/home.png");
  } else {
    addNewMenuButton("", 2, ":/img/img/home.png");
  }
  addNewMenuButton("", 6, ":/img/img/add.png");
  addNewMenuButton("", 6, ":/img/img/add.png");
  addNewMenuButton("", 6, ":/img/img/add.png");
  addNewMenuButton("", 6, ":/img/img/add.png");
  addNewMenuButton("", 6, ":/img/img/add.png");
  addNewMenuButton("", 6, ":/img/img/add.png");
  addNewMenuButton("", 3, ":/img/img/add.png");
}

void MainWindow::addNewMenuButton(QString text, int type, QString iconPath) {
  MenuButton *newMenuButton = new MenuButton(this, text, type, iconPath);
  QHBoxLayout *newMenuButtonLayout = new QHBoxLayout;
  newMenuButtonLayout->setAlignment(Qt::AlignHCenter);
  if (type == 2 || type == 3 || type == 4) {
    if (type == 3) {
      this->menuPanelPage_Layout->addStretch();
    }
    MenuButton *newMenuArrowLeft =
        new MenuButton(this, "", -51, ":/img/img/leftArrow.png");
    MenuButton *newMenuArrowRight = new MenuButton(this, "", -52, "");
    if (type == 2 || type == 4) {
      this->menuArrowTop = newMenuArrowRight;
    } else if (type == 3) {
      this->menuArrowBot = newMenuArrowRight;
    }
    newMenuButtonLayout->addWidget(newMenuArrowLeft);
    newMenuButtonLayout->addWidget(newMenuButton);
    newMenuButtonLayout->addWidget(newMenuArrowRight);
    this->addButton = newMenuButton;
    this->menuPanelPage_Layout->addLayout(newMenuButtonLayout);
    if (type != 3) {
      this->menuPanelPage_Layout->addStretch();
      newMenuButton->indexOfContentPage = 0;
      if (this->amountOfContentPages == 0) {
        addContentPage(newMenuButton);
      }
    }
    if (!this->menuPanel->currentIndex()) {
      newMenuArrowLeft->setIcon(QPixmap(""));
    }
  } else {
    newMenuButtonLayout->addWidget(newMenuButton);
    this->menuPanelPage_Layout->addLayout(newMenuButtonLayout);
  }
}

void MainWindow::addButtonWasClicked(MenuButton *button) {
  if (button->type == 3) {
    QWidget *menuPanelPage = new QWidget;
    this->menuPanel->addWidget(menuPanelPage);
    QVBoxLayout *menuPanelPage_Layout = new QVBoxLayout;
    menuPanelPage_Layout->setSpacing(35);
    menuPanelPage->setLayout(menuPanelPage_Layout);
    this->menuPanelPage_Layout = menuPanelPage_Layout;
    this->menuPanel->setCurrentIndex(this->menuPanel->currentIndex() + 1);
    this->menuArrowTop->setIcon(QPixmap(":/img/img/rightArrow.png"));
    this->menuArrowBot->setIcon(QPixmap(":/img/img/rightArrow.png"));
    this->addButton->setIcon(QPixmap(""));
    this->addButton->disconnect(this->addButton,
                                SIGNAL(addButtonWasClicked(MenuButton *)), this,
                                SLOT(addButtonWasClicked(MenuButton *)));
    fillMenuPanel();
  } else if (button->type == 6) {
    button->type = 0;
    button->connectMenuButton(this);
    addContentPage(button);
  }
}

void MainWindow::addContentPage(MenuButton *button) {
  button->indexOfContentPage = this->amountOfContentPages++;
  QWidget *contentPanelPage = new QWidget;
  QGridLayout *contentPanelPage_Layout = new QGridLayout;
  contentPanelPage_Layout->setAlignment(Qt::AlignCenter);
  contentPanelPage->setLayout(contentPanelPage_Layout);
  this->contentPanel->addWidget(contentPanelPage);
  int x = 0, y = 0;
  for (int i = 0; i <= button->indexOfContentPage; i++) {
    QPushButton *count = new QPushButton;
    QString number, label;
    if (!button->indexOfContentPage) {
      count->setStyleSheet("background-color: rgb(60,0,100); color: white;"
                           "border-radius: 75px");
      label = "ContentPage_Home";
      count->setText(label);
      count->setFont(QFont("Century Gothic", 10, QFont::Normal));
      count->setFixedSize(150, 150);
    } else if (button->indexOfContentPage < 16) {
      count->setStyleSheet("background-color: rgb(60,60,60); color: white;"
                           "border-radius: 75px");
      number.setNum(button->indexOfContentPage);
      label = "ContentPage_" + number;
      count->setText(label);
      count->setFont(QFont("Century Gothic", 10, QFont::Normal));
      count->setFixedSize(150, 150);
    } else {
      count->setStyleSheet("background-color: rgb(60,0,100); color: white;"
                           "border-radius: 37px");
      number.setNum(button->indexOfContentPage);
      label = "ContentPage_" + number;
      count->setText(label);
      count->setFont(QFont("Century Gothic", 5, QFont::Normal));
      count->setFixedSize(74, 74);
    }
    if (i >= button->indexOfContentPage && i) {
      break;
    }
    contentPanelPage_Layout->addWidget(count, x, y++);
    if (y == 5 && button->indexOfContentPage < 16) {
      y = 0;
      x++;
    } else if (y == 7) {
      y = 0;
      x++;
    }
  }
  QString number, label;
  number.setNum(button->indexOfContentPage);
  if (button->indexOfContentPage) {
    button->setFixedSize(184, 40);
    button->setIcon(QPixmap(""));
    label = "MenuButton_" + number;
    button->setText(label);
  }
}

void MainWindow::emitMenuButtonWasClicked() { emit menuButtonWasClicked(); }

void MainWindow::minimizeApp() {
  this->showMinimized();
  emit systemButtonWasClicked();
}

void MainWindow::maximizeApp() {
  if (isMaximized()) {
    this->showNormal();
    emit systemButtonWasClicked();
  } else {
    this->showMaximized();
  }
}

void MainWindow::closeApp() {
  emit systemButtonWasClicked();
  this->close();
}

void MainWindow::setContentPanelIndex(int indexOfContentPage) {
  this->contentPanel->setCurrentIndex(indexOfContentPage);
}

void MainWindow::setMenuPanelIndex(int index) {
  if (index == -51 && this->menuPanel->currentIndex()) {
    this->menuPanel->setCurrentIndex(this->menuPanel->currentIndex() - 1);
  } else if (index == -52 &&
             this->menuPanel->currentIndex() + 1 < this->amountOfMenuPages) {
    setMenuPanelIndex(this->menuPanel->currentIndex() + 1);
  } else if (index >= 0) {
    this->menuPanel->setCurrentIndex(index);
  }
}
