job "job1" "POST" {
    query = {
        token = ""
    }
    url = ""
}

pipeline "pipeline1" {
    job "first job" "job1" {}
    job "second job" "job1" {
        requires = ["first job"]
    }
}

pipeline "pipeline2" {
    job "first job" "job1" {}
}
