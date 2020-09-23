(ns clj-journal.core-test
  (:require [clojure.test :refer :all]
            [clj-journal.core :refer :all]))

;; TODO: write a macro to generate '(is ...)' forms to compare all the files in a pair of directories
(defmacro cmp-dir-files [dir1-f dir2-f]
  (loop [files1 (file-seq (eval dir1-f)) files2 (file-seq (eval dir2-f)) forms '()]
    (if
        (and (empty? files1) (empty? files2))
      forms
      (if
          (or (.isDirectory (first files1)) (.isDirectory (first files2)))
        (do
          (println
           (str
            "Warning: skipping sub-directory in directory diff: "
            (.getPath (first files1))
            " (in directory 1), "
            (.getPath (first files2))
            " (in directory 2)"))
          (recur (rest files1) (rest files2) forms))
        (recur
         (rest files1)
         (rest files2)
         (cons
          `(is (= ~(slurp (first files1)) ~(slurp (first files2)))) forms))))))

(defn copy-test-resource [file-name target-dir]
  (clojure.java.io/copy
   (clojure.java.io/file
    (clojure.java.io/resource (str "test/" file-name)))
   (clojure.java.io/file (str target-dir "/" file-name))))

(def ^:dynamic tmp-dir "")
(use-fixtures
  :each
  (fn [f]
    (binding
     [tmp-dir
      (str "/tmp/clj-journal/test/" (.getTime (java.util.Date.)))]
      (println (str "Using tmp dir: " tmp-dir))
      (.mkdirs (clojure.java.io/file tmp-dir))
      (.mkdir (clojure.java.io/file (str tmp-dir "/may")))
      (copy-test-resource "may/may01.note" tmp-dir)
      (copy-test-resource "may/may02.note" tmp-dir)
      (copy-test-resource "may/may11.note" tmp-dir)
      (copy-test-resource "may/may21.note" tmp-dir)
      (f))))

(deftest migrate-day-test
  (testing
   "migrate-day"
    (let
     [inst (java.util.Calendar/getInstance)]
      (.set inst (java.util.Calendar/YEAR) 2020)
      (.set inst (java.util.Calendar/MONTH) 4)
      (.set inst (java.util.Calendar/DAY_OF_MONTH) 25)
      (migrate-day tmp-dir inst)
      (is
       (=
        (slurp
         (clojure.java.io/file
          (clojure.java.io/resource "test/expected-day.note")))
        (slurp (clojure.java.io/file (str tmp-dir "/may/may25.note"))))))))

(deftest migrate-month-test
  (testing
   "migrate-month"
    (let
        [inst
         (java.util.Calendar/getInstance)
         dir-f
         (clojure.java.io/file (str tmp-dir "/jun"))
         expected-dir-f
         (clojure.java.io/file (clojure.java.io/resource "test/expected-jun"))]
      (.set inst (java.util.Calendar/YEAR) 2020)
      (.set inst (java.util.Calendar/MONTH) 5)
      (.set inst (java.util.Calendar/DAY_OF_MONTH) 1)
      (migrate-month tmp-dir inst)
      (is
       (=
        (slurp
         (clojure.java.io/file
          (clojure.java.io/resource "test/expected-jun/tasks.note")))
        (slurp (clojure.java.io/file (str tmp-dir "/jun/tasks.note")))))
      (is
       (=
        (slurp
         (clojure.java.io/file
          (clojure.java.io/resource "test/expected-jun/calendar.txt")))
        (slurp (clojure.java.io/file (str tmp-dir "/jun/calendar.txt")))))
      ;;(cmp-dir-files expected-dir-f dir-f)
      )))
