(ns clj-journal.stream)

(defn mapper
  [reader callback]
  (loop [lines (line-seq (java.io.BufferedReader. reader))]
    (let [line (first lines)]
      (if line (do
                 (callback line)
                 (recur (rest lines)))
          ))))
