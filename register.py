# -*- coding: utf-8 -*-
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import Select
from selenium.common.exceptions import NoSuchElementException
from selenium.common.exceptions import NoAlertPresentException
import unittest, time, re

class Register(unittest.TestCase):
    def setUp(self):
        self.driver = webdriver.Firefox()
        self.driver.implicitly_wait(30)
        self.base_url = "https://id.asics.theplant-dev.com/"
        self.verificationErrors = []
        self.accept_next_alert = True
    
    def test_register(self):
        driver = self.driver
        driver.get(self.base_url + "/app?locale=en-US")
        driver.find_element_by_link_text("Register").click()
        driver.find_element_by_name("email").clear()
        driver.find_element_by_name("email").send_keys("azuma+1@theplant.jp")
        driver.find_element_by_name("password").clear()
        driver.find_element_by_name("password").send_keys("123456")
        driver.find_element_by_name("confirmed_password").clear()
        driver.find_element_by_name("confirmed_password").send_keys("123456")
        Select(driver.find_element_by_name("country")).select_by_visible_text("China")
        Select(driver.find_element_by_name("day")).select_by_visible_text("16")
        Select(driver.find_element_by_name("month")).select_by_visible_text("12")
        Select(driver.find_element_by_name("year")).select_by_visible_text("1916")
        driver.find_element_by_css_selector("input[type=\"checkbox\"]").click()
        driver.find_element_by_xpath("//button[@type='submit']").click()
        try: self.assertEqual("New account created", driver.find_element_by_xpath("//div[@id='mount']/div/div/div[2]/section/h2/span[2]").text)
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
        self.driver.quit()
        self.assertEqual([], self.verificationErrors)

if __name__ == "__main__":
    unittest.main()
