
## generate runable test
./bang.py ./checkprodution.py


## browser config
https://wiki.saucelabs.com/display/DOCS/Platform+Configurator#/



## note
java -jar selenium-server-standalone-2.53.0.jar -role hub
java -jar selenium-server-standalone-2.53.0.jar -role node  -hub http://localhost:4444/grid/register -nodeConfig nodeConfig.json