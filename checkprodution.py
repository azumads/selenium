# -*- coding: utf-8 -*-
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import Select
from selenium.common.exceptions import NoSuchElementException
from selenium.common.exceptions import NoAlertPresentException
import unittest, time, re

class Checkproduction(unittest.TestCase):
    def setUp(self):
        self.driver = webdriver.Firefox()
        self.driver.implicitly_wait(30)
        self.base_url = "http://demo.getqor.com/"
        self.verificationErrors = []
        self.accept_next_alert = True
    
    def test_checkproduction(self):
        driver = self.driver
        driver.get(self.base_url + "/admin")
        driver.find_element_by_name("email").clear()
        driver.find_element_by_name("email").send_keys("dev@getqor.com")
        driver.find_element_by_name("password").clear()
        driver.find_element_by_name("password").send_keys("testing")
        driver.find_element_by_css_selector("button.button.button__primary").click()
        driver.get(self.base_url + "/admin")
        driver.find_element_by_link_text("Products").click()
        try: self.assertEqual("QOR Jacket", driver.find_element_by_xpath("//main[@id='content']/div[2]/div/table/tbody/tr[10]/td[2]/div").text)
        except AssertionError as e: self.verificationErrors.append(str(e))
        driver.find_element_by_css_selector("a.mdl-navigation__link > i.material-icons").click()
    
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
        self.driver.quit()
        self.assertEqual([], self.verificationErrors)

if __name__ == "__main__":
    unittest.main()
