namespace :app do
  task :install do
    sh "git pull && go build"
	sh "cp pushups /usr/local/bin/pushups/"
	sh "cp -R static/ /usr/local/bin/pushups/"
  end
end