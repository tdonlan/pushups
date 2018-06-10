namespace :app do
  task :install do
    sh "sudo service pushups stop"
    sh "git pull && go build"
    sh "sudo cp pushups /usr/local/bin/pushups/"
    sh "sudo cp -R static/ /usr/local/bin/pushups/"
    sh "sudo service pushups start"
  end
end