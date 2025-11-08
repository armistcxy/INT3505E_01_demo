BASE_URL="http://localhost:8080"

N=5000
shuf -i 1-300 -n "$N" -r | awk -v base="$BASE_URL" '{print "GET " base "/api/v1/books/" $1}' > targets.txt

head -n 5 targets.txt