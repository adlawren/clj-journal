(ns clj-journal.core-test
  (:require [clojure.test :refer :all]
            [clj-journal.core :refer :all]))

(defn copy-dir [src-f dest-f]
  (.mkdir dest-f)
  (loop [src-files (.listFiles src-f)]
    (if
        (not (empty? src-files))
      (let [f (first src-files)]
        (if
            (not (= (.getPath f) (.getPath src-f)))
          (let
              [target-f
               (clojure.java.io/file (str (.getPath dest-f) "/" (.getName f)))]
            (if
                (.isDirectory f)
              (copy-dir f target-f)
              (clojure.java.io/copy f target-f))))
        (recur (rest src-files))))))

(defn sorted-files [dir-f]
  (sort
   (fn [f1 f2] (compare (.getPath f1) (.getPath f2)))
   (.listFiles dir-f)))

(defn cmp-dirs [dir1-f dir2-f]
  (loop [files1
         (sorted-files dir1-f)
         files2
         (sorted-files dir2-f)]
    (if
        (or (not (empty? files1)) (not (empty? files2)))
      (let [file1 (first files1) file2 (first files2)]
        (if
            (and
             (= (.getPath dir1-f) (.getPath file1))
             (= (.getPath dir2-f) (.getPath file2)))
          (recur
           (rest files1)
           (rest files2))
          (do
            (if
                (or (.isDirectory file1) (.isDirectory file2))
              (cmp-dirs file1 file2)
              (clojure.test/is
               (= (slurp file1) (slurp file2))
               (str (.getPath file1) " does not match " (.getPath file2))))
            (recur
             (rest files1)
             (rest files2))))))))

(def ^:dynamic tmp-dir "")
(use-fixtures
  :each
  (fn [f]
    (binding
        [tmp-dir
         (str "/tmp/clj-journal/test/" (.getTime (java.util.Date.)))]
      (println (str "Using tmp dir: " tmp-dir))
      (.mkdirs (clojure.java.io/file (str tmp-dir "/2019")))
      (copy-dir
       (clojure.java.io/file (clojure.java.io/resource "test/dec"))
       (clojure.java.io/file (str tmp-dir "/2019/dec")))
      (f))))

(deftest migrate-day-test
  (testing
      "migrate-day"
    (let
        [cal-inst (java.util.Calendar/getInstance)]
      (.set cal-inst (java.util.Calendar/YEAR) 2019)
      (.set cal-inst (java.util.Calendar/MONTH) 11)
      (.set cal-inst (java.util.Calendar/DAY_OF_MONTH) 25)
      (migrate-day tmp-dir cal-inst)
      (cmp-dirs
       (clojure.java.io/file (clojure.java.io/resource "test/expected-dec"))
       (clojure.java.io/file (str tmp-dir "/2019/dec"))))))

(deftest migrate-month-test
  (testing
      "migrate-month"
    (let
        [cal-inst (java.util.Calendar/getInstance)]
      (.set cal-inst (java.util.Calendar/YEAR) 2020)
      (.set cal-inst (java.util.Calendar/MONTH) 0)
      (.set cal-inst (java.util.Calendar/DAY_OF_MONTH) 1)
      (migrate-month tmp-dir cal-inst)
      (cmp-dirs
       (clojure.java.io/file (clojure.java.io/resource "test/expected-jan"))
       (clojure.java.io/file (str tmp-dir "/2020/jan"))))))
