#!/usr/bin/env python
# -*- coding: utf-8 -*-
import sys
import os
import re
import csv
import stat


imports = '''from sauceclient import SauceClient
import new, sys, csv\n
'''

onbrowsers = '''browsers = [{"platform": "Mac OS X 10.9",
             "browserName": "chrome",
             "version": "31",
             "screenResolution":"1280x1024"}]'''
onplatforms = '''
def on_platforms(platforms):
    def decorator(base_class):
        module = sys.modules[base_class.__module__].__dict__
        for i, platform in enumerate(platforms):
            d = dict(base_class.__dict__)
            d['desired_capabilities'] = platform
            name = "%s_%s" % (base_class.__name__, i + 1)
            module[name] = new.classobj(name, (base_class,), d)
    return decorator

@on_platforms(browsers)\n'''

setups = '''        self.desired_capabilities['name'] = self.id()
        sauce_url = "http://%s:%s@ondemand.saucelabs.com:80/wd/hub"
        self.driver = webdriver.Remote(
            command_executor=sauce_url % ("testsaucelaber", "097cc55a-4c6e-4ee7-bdb9-0868ecb01b72"),
            desired_capabilities=self.desired_capabilities)
        self.sauce_client = SauceClient("testsaucelaber", "097cc55a-4c6e-4ee7-bdb9-0868ecb01b72")            
        self.driver.implicitly_wait(10)
'''

teardowns = '''        print("Link to your job: https://saucelabs.com/jobs/%s" % self.driver.session_id)
        try:
            if sys.exc_info() == (None, None, None):
                self.sauce_client.jobs.update_job(self.driver.session_id, passed=True)
            else:
                self.sauce_client.jobs.update_job(self.driver.session_id, passed=False)
        finally:
            self.driver.quit()\n'''
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
                    newfile.write(onbrowsers)
                    newfile.write(onplatforms)
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
                elif val.startswith("    def tearDown(self)"):
                    newfile.write(val)
                    newfile.write(teardowns)
                    skip = 2
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
