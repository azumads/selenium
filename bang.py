#!/usr/bin/env python
# -*- coding: utf-8 -*-
import sys
import os
import re
import csv
import stat


imports = '''import new, sys, csv\n'''


setups = '''        self.driver = webdriver.Chrome()
        self.driver.set_window_size(1280, 1024)          
        self.driver.implicitly_wait(15)
'''

sendkeyRe = "send_keys\\(\"(.+?)\""
selectRe = "select_by_visible_text\\(\"(.+?)\""
assertEqualRe = "assertEqual\\(\"(.+?)\","

def bang():
    args = sys.argv

    if len(args) < 2:
        print "wrong arguments"
        print args
        return
    name = os.path.basename(args[1])
    directory = os.path.dirname(args[1])
    csvFile = ""
    ProvideCsv = False
    if len(args) == 3:
        csvFile  = args[2]

    if csvFile != "":
        if os.path.exists(csvFile) != True:
            print "csv file don't exists"
            return
        else:
            ProvideCsv = True

    if csvFile == "":
       csvFile =  directory+"/"+name[0:len(name)-3]+".csv"

    # print name
    # print directory
    with open("browserConfig", 'r') as f:
        browsers = f.read()
        if browsers != "":
            onbrowsers = browsers

    skip = 0
    inTest = False
    csvIndex = 0
    row = []
    with open(directory+"/B"+name, 'w') as newfile:
        with open(args[1], 'r') as originalfile:
            for line, val in enumerate(originalfile.readlines()):
                if skip != 0:
                    skip = skip - 1
                    continue
                if line == 0 and val != "#!/usr/bin/env python":
                    newfile.write('#!/usr/bin/env python\n')
                    newfile.write(val)
                elif val.startswith("class"):
                    newfile.write(imports)
                #     newfile.write(onbrowsers)
                #     newfile.write(onplatforms)
                    newfile.write(val)
                elif val.startswith("    def setUp("):
                    newfile.write(val)
                    newfile.write(setups)
                    skip = 2
                elif val.startswith("        driver = self.driver"):
                    newfile.write(val)
                    newfile.write("        with open(\"" + csvFile + "\", 'r') as csvfile:\n")
                    newfile.write("            reader = csv.reader(csvfile)\n")
                    newfile.write("            for data in reader:\n") 
                    inTest = True
                elif val.startswith("    def is_element_present"):
                    inTest = False
                    newfile.write(val)
                # elif val.startswith("    def tearDown(self)"):
                #     newfile.write(val)
                #     newfile.write(teardowns)
                #     skip = 2
                elif inTest and val.find("send_keys(")>0:
                    # print re.search(sendkeyRe, val).group(1)
                    newfile.write("        ") 
                    newfile.write(re.sub(sendkeyRe,"send_keys(data["+str(csvIndex)+"]",val))
                    row.append(re.search(sendkeyRe, val).group(1))
                    csvIndex += 1
                elif inTest and val.find("select_by_visible_text(")>0:
                    # print re.search(selectRe, val).group(1)
                    newfile.write("        ") 
                    newfile.write(re.sub(selectRe,"select_by_visible_text(data["+str(csvIndex)+"]",val))
                    row.append(re.search(selectRe, val).group(1))
                    csvIndex += 1
                elif inTest and val.find("assertEqual(")>0:
                    # print re.search(assertEqualRe, val).group(1)
                    newfile.write("        ") 
                    newfile.write(re.sub(assertEqualRe,"assertEqual(data["+str(csvIndex)+"],",val))
                    row.append(re.search(assertEqualRe, val).group(1))
                    csvIndex += 1
                else:
                    if inTest :
                        newfile.write("        ") 
                    newfile.write(val) 
    if not ProvideCsv:
        with open(csvFile, 'w') as csvfile:
            writer = csv.writer(csvfile)
            writer.writerow(row)


    os.chmod(directory+"/B"+name,stat.S_IRWXU)
    print directory+"/B"+name


if __name__=='__main__':
    bang()
