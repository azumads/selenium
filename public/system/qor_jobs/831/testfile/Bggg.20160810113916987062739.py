#!/usr/bin/env python
# -*- coding: utf-8 -*-
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import Select
from selenium.common.exceptions import NoSuchElementException
from selenium.common.exceptions import NoAlertPresentException
import unittest, time, re

from sauceclient import SauceClient
import new, sys, csv

browsers = [{"platform": "OS X 10.9",
             "browserName": "chrome",
             "version": "31",
             "screenResolution":"1280x1024"},
            {"platform": "Windows 8.1",
             "browserName": "internet explorer",
             "version": "11",
             "screenResolution":"1280x1024"},
             {"platform": "OS X 10.11",
             "browserName": "safari",
             "version": "9.0",
             "screenResolution":"2048x1536"}]

def on_platforms(platforms):
    def decorator(base_class):
        module = sys.modules[base_class.__module__].__dict__
        for i, platform in enumerate(platforms):
            d = dict(base_class.__dict__)
            d['desired_capabilities'] = platform
            name = "%s_%s" % (base_class.__name__, i + 1)
            module[name] = new.classobj(name, (base_class,), d)
    return decorator

@on_platforms(browsers)
class Ggg(unittest.TestCase):
    def setUp(self):
        self.desired_capabilities['name'] = self.id()
        sauce_url = "http://%s:%s@ondemand.saucelabs.com:80/wd/hub"
        self.driver = webdriver.Remote(
            command_executor=sauce_url % ("testsaucelaber", "097cc55a-4c6e-4ee7-bdb9-0868ecb01b72"),
            desired_capabilities=self.desired_capabilities)
        self.sauce_client = SauceClient("testsaucelaber", "097cc55a-4c6e-4ee7-bdb9-0868ecb01b72")            
        self.driver.implicitly_wait(10)
        self.base_url = "http://demo.getqor.com/"
        self.verificationErrors = []
        self.accept_next_alert = True
    
    def test_ggg(self):
        driver = self.driver
        with open("public/system/qor_jobs/831/testfile/ggg.20160810113916987062739.csv", 'r') as csvfile:
            reader = csv.reader(csvfile)
            for data in reader:
                driver.get(self.base_url + "/auth/login")
                driver.find_element_by_name("email").clear()
                driver.find_element_by_name("email").send_keys(data[0])
                driver.find_element_by_name("password").clear()
                driver.find_element_by_name("password").send_keys(data[1])
                driver.find_element_by_css_selector("button.button.button__primary").click()
                driver.find_element_by_link_text("Admin Dashboard").click()
                driver.find_element_by_link_text("Products").click()
                try: self.assertEqual(data[2], driver.find_element_by_xpath("//main[@id='content']/div[2]/div/table/tbody/tr/td[2]/div").text)
                except AssertionError as e: self.verificationErrors.append(str(e))
            
    def is_element_present(self, how, what):
        try: self.driver.find_element(by=how, value=what)
        except NoSuchElementException as e: return False
        return True
    
    def is_alert_present(self):
        try: self.driver.switch_to_alert()
        except NoAlertPresentException as e: return False
        return True
    
    def close_alert_and_get_its_text(self):
        try:
            alert = self.driver.switch_to_alert()
            alert_text = alert.text
            if self.accept_next_alert:
                alert.accept()
            else:
                alert.dismiss()
            return alert_text
        finally: self.accept_next_alert = True
    
    def tearDown(self):
        print("Link to your job: https://saucelabs.com/jobs/%s" % self.driver.session_id)
        try:
            if sys.exc_info() == (None, None, None):
                self.sauce_client.jobs.update_job(self.driver.session_id, passed=True)
            else:
                self.sauce_client.jobs.update_job(self.driver.session_id, passed=False)
        finally:
            self.driver.quit()

if __name__ == "__main__":
    unittest.main()
