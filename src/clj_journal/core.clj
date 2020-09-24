(ns clj-journal.core
  (:require
   [clj-journal.stream :as stream]
   [clojure.tools.cli])
  (:gen-class))

(defn month-name [calendar-inst]
  (clojure.string/lower-case
   (.getDisplayName
    calendar-inst
    (java.util.Calendar/MONTH)
    (java.util.Calendar/LONG_STANDALONE)
    (java.util.Locale. "en-US"))))

(defn day [calendar-inst]
  (.get calendar-inst (java.util.Calendar/DAY_OF_MONTH)))

(defn days-in-month [calendar-inst]
  (.getActualMaximum calendar-inst (java.util.Calendar/DAY_OF_MONTH)))

(defn day-name [calendar-inst]
  (.getDisplayName
   calendar-inst
   (java.util.Calendar/DAY_OF_WEEK)
   (java.util.Calendar/LONG_STANDALONE)
   (java.util.Locale. "en-US")))

(defn note-files [path]
  (filter
   (fn [f] (not (nil? (re-find #".note$" (.getPath f)))))
   (file-seq (clojure.java.io/file path))))

(defprotocol NoteTreeP (serialize [this]))
(defrecord NoteTree [note children]
  NoteTreeP
  (serialize [this]
    (str
     (:text note)
     \newline
     (clojure.string/join (map (fn [nt] (serialize nt)) children)))))

(def bullet-regex-str "[?\\-*x~><]")
(def leading-whitespace-regex-str "^[\\s]*")
(def trailing-whitespace-regex-str "[\\s]+")
(def note-regex
  (re-pattern
   (str
    leading-whitespace-regex-str
    bullet-regex-str
    trailing-whitespace-regex-str)))
(defprotocol NoteP
  (isBullet? [this])
  (isUncompletedTask? [this])
  (depth [this]))
(defrecord Note [text]
  NoteP
  (isBullet? [this] (not (nil? (re-find note-regex text))))
  (isUncompletedTask? [this]
    (let [bullet (re-find note-regex text)]
      (if (nil? bullet) false (clojure.string/includes? bullet "*"))))
  (depth [this]
    (let [bullet (re-find note-regex text)]
      (if
          (nil? bullet)
        0
        (let [whitespace
              (take
               (-
                (count
                 (re-find
                  (re-pattern
                   (str leading-whitespace-regex-str bullet-regex-str))
                  bullet))
                1) bullet)]
          (if
              (> (count (filter (fn [c] (not (= c \space))) whitespace)) 0)
            (println
             (str
              "Error: can't determine depth because line contains mixed whitespace: \""
              text
              "\""))
            (count whitespace)))))))

(defn parse-notes [lines] (map (fn [l] (Note. l)) lines))

(defn filter-bullets [notes]
  (if
      (empty? notes)
    notes
    (apply
     vector
     (reverse
      (loop [prev-text (:text (first notes)) new-notes '() notes (rest notes)]
        (if
            (empty? notes)
          (cons (Note. prev-text) new-notes)
          (let [text (:text (first notes))]
            (if
                (.isBullet? (first notes))
              (recur text (cons (Note. prev-text) new-notes) (rest notes))
              (recur (str prev-text \newline text) new-notes (rest notes))))))))))

(def ^:dynamic current-note-idx 0)
(defn parse-note-trees [notes]
  (binding [current-note-idx 0]
    ((fn parse-note-trees-internal [notes current-depth]
       (if
           (or
            (>= current-note-idx (count notes))
            (< (.depth (.get notes current-note-idx)) current-depth))
         '()
         (let [idx current-note-idx]
           (set! current-note-idx (+ current-note-idx 1))
           (cons
            (NoteTree.
             (.get notes idx)
             (if
                 (= idx (- (count notes) 1))
               '()
               (let [next-depth (.depth (.get notes (+ idx 1)))]
                 (if
                     (> next-depth current-depth)
                   (parse-note-trees-internal notes next-depth)
                   '()))))
            (parse-note-trees-internal notes current-depth))))) notes 0)))

(defn migrate-note-trees [note-trees]
  (if
      (empty? note-trees)
    '()
    (if
        (.isUncompletedTask? (:note (first note-trees)))
      (cons (first note-trees) (migrate-note-trees (rest note-trees)))
      (let [migrated-children
            (migrate-note-trees (:children (first note-trees)))]
        (if
            (empty? migrated-children)
          (migrate-note-trees (rest note-trees))
          (cons
           (NoteTree. (:note (first note-trees)) migrated-children)
           (migrate-note-trees (rest note-trees))))))))

;; TODO: try to leverage the NoteTree data structure instead?
(defn update-uncompleted-task-lines [note-files]
  (loop [note-files note-files]
    (if (not (empty? note-files))
      (do
        (let [note-file (first note-files)]
          (spit
           note-file
           (str
            (clojure.string/join
             \newline
             (map
              (fn [l]
                (clojure.string/replace
                 l
                 (re-pattern
                  (str leading-whitespace-regex-str "\\*" trailing-whitespace-regex-str))
                 (fn [s] (clojure.string/replace s "*" ">"))))
              (clojure.string/split (slurp note-file) #"\n"))) \newline)))
        (recur (rest note-files))))))

;; TODO: update migrated files
(defn migrate [note-files target-file]
  (spit
   target-file
   (clojure.string/join
    ""
    (map
     (fn [nt] (.serialize nt))
     (migrate-note-trees
      (flatten
       (map
        (fn [f]
          (parse-note-trees
           (filter-bullets
            (parse-notes (clojure.string/split (slurp f) #"\n")))))
        note-files))))))
  (update-uncompleted-task-lines note-files))

(defn migrate-day [path calendar-inst]
  (let [month
        (month-name calendar-inst)
        day
        (day calendar-inst)
        target-file
        (clojure.java.io/file
         (str path "/" month "/" month day ".note"))]
    (if
        (.exists target-file)
      (println (str "Error: " (.getPath target-file) " already exists"))
      (migrate
       (filter
        (fn [f]
          (not
           (nil?
            (re-find
             (re-pattern (str month "[0-9]+.note$"))
             (.getPath f)))))
        (note-files (str path "/" month)))
       target-file))))

(defn create-calendar-file [path calendar-inst]
  (spit
   (clojure.java.io/file
    (str path "/" (month-name calendar-inst) "/calendar.txt"))
   (str
    (clojure.string/join
     \newline
     (map
      (fn [idx]
        (let
            [inst (java.util.Calendar/getInstance)]
          (.set inst (java.util.Calendar/DAY_OF_MONTH) idx)
          (str (day-name inst) " " (+ idx 1))))
      (take
       (days-in-month (java.util.Calendar/getInstance))
       (iterate (fn [x] (+ x 1)) 0))))
    \newline)))

(defn migrate-month [path calendar-inst]
  (let [prev-month-inst
        (let
            [new-inst (.clone calendar-inst)]
          (.add new-inst (java.util.Calendar/MONTH) -1)
          new-inst)
        prev-month
        (month-name prev-month-inst)
        current-month
        (month-name calendar-inst)
        month-dir
        (clojure.java.io/file (str path "/" current-month))
        target-file
        (clojure.java.io/file
         (str path "/" current-month "/tasks.note"))]
    (if
        (.exists month-dir)
      (println (str "Error: " (.getPath month-dir) " already exists"))
      (do
        (.mkdir month-dir)
        (create-calendar-file path calendar-inst)
        (migrate (note-files (str path "/" prev-month)) target-file)))))

(def cli-options
  [["-m" "--migrate-day"] ["-M" "--migrate-month"] ["-h" "--help"]])

;; TODO
;; (System/getProperty "user.dir") ;; cwd
(def target-dir (str (System/getenv "HOME") "/test-dir"))

(defn -main [& args]
  (let [opts (clojure.tools.cli/parse-opts args cli-options)]
    (cond
      (:help opts)
      (println "Use '-m' or '--migrate-day' to run a daily migration, '-M' or '--migrate-month' to run a monthly migration")
      (:migrate-day opts)
      (migrate-day target-dir (java.util.Calendar/getInstance))
      (:migrate-month opts)
      (migrate-month target-dir (java.util.Calendar/getInstance)))))