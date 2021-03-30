# How to run the test

## 1. use the selenium IDE(firefox add-ons) to record test case

	selenium IDE install url: https://addons.mozilla.org/en-US/firefox/addon/selenium-ide/

## 2. export test case as python2/ unittest / webdriver. (file name end with .py)


![ScreenShot](https://raw.githubusercontent.com/azumads/selenium/master/export.png)


## 3. use the export test file to create and run test case in admin 

 	admih url: http://192.168.1.203:8000/admin






## generate runable test
./bang.py ./checkprodution.py


## browser config
https://wiki.saucelabs.com/display/DOCS/Platform+Configurator#/



## note
java -jar selenium-server-standalone-2.53.0.jar -role hub
java -jar selenium-server-standalone-2.53.0.jar -role node  -hub http://localhost:4444/grid/register -nodeConfig nodeConfig.json