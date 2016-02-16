require 'bundler/setup'
require 'shellwords'

@golang_version = "1.5.3"
@app_name = "kanban-stats"
@blue_stack = "blue"
@green_stack = "green"

@current_dir = File.dirname(__FILE__)
@deploy_directory = "#{@current_dir}/deploy"
@build_directory = "#{@current_dir}/build"
@version_file = "#{@build_directory}/version"

@pod_timeout_secs = 300
@smoketest_timeout_secs = 30

@service_port = 80
@container_port = 8080

def validate_golang_version
	#validate go version
  version = `go version`.chomp
  raise "Golang is not the correct version, expected #{@golang_version}, got '#{version}'" unless version.include?(@golang_version)
end

def get_commit_hash
	# get version
  commit = `git rev-parse --short HEAD`.chomp
  local_changes = `git status --porcelain`
  changes = ""
  # changes = "+CHANGES" if (local_changes != "")
  "#{commit}#{changes}"
end

def get_os
  if (/darwin/ =~ RUBY_PLATFORM) != nil
    return "darwin"
  else
    return "linux"
  end
end

def create_version_file(version)
	File.write("#{@version_file}", "{ \"version\": \"#{version}\" }")
end

task :build_app do
	validate_golang_version
	commit_note = get_commit_hash

	create_directory(@build_directory)
	output_file = Shellwords.escape("#{@build_directory}/#{@app_name}")

	# Restore dependencies
	sh "go get github.com/tools/godep"
	sh "godep restore"

	# build target
	sh "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -o #{output_file} -ldflags \"-X main.GitCommit=#{commit_note}\""
	# build local for non-linux OS
	os = get_os
	if os != "linux"
		output_file = "#{output_file}_local"
		sh "CGO_ENABLED=0 godep go build -o #{output_file} -ldflags \"-X main.GitCommit=#{commit_note}\""
	end

	version = `#{output_file} -version`.strip
	create_version_file(version)
end

task :unit_test do
	# Restore dependencies
	sh "go get github.com/tools/godep"
	sh "godep restore"

	sh "godep go test -v -timeout=5s ./..."
end

task :build_image do
	create_directory(@build_directory)

	# Copy Dockerfile to build_directory(for developer workstations)
	file = "Dockerfile"
	if File.exist?(file)
		FileUtils.cp(file, @build_directory)
	end

	# make dure file is executable before copying to image
	File.chmod(0744, "#{@build_directory}/#{@app_name}")

	sh "docker build --no-cache --build-arg APP_NAME=#{@app_name} -t #{@app_name} #{Shellwords.escape(@build_directory)}"
end

task :release_image do
	account = docker_account
	version = get_version

	sh "docker tag -f #{@app_name} #{account}/#{@app_name}:#{version}"
	sh "docker push #{account}/#{@app_name}:#{version}"
end

task :build => [:build_app, :unit_test, :build_image, :release_image]

def get_current_context
	`kubectl config view -o jsonpath='{.current-context}'`
end

def create_directory(path)
  Dir.mkdir(path) unless Dir.exist?(path)
end

def render_templates(vars)
	create_directory(@deploy_directory)

	require 'erb'
	Dir.glob("#{@current_dir}/templates/*.yaml.erb") do |file|
		renderer = ERB.new(IO.read(file))
		output = renderer.result(Namespace.new(vars).get_binding)

		file_name = File.basename(file, ".*")
		output_file = File.join(@deploy_directory, "/#{file_name}")
		IO.write(output_file, output)
	end
end

def dns_domain
	context = get_current_context
	api_server_url = `kubectl config view -o jsonpath='{.clusters[?(@.name==\"#{context}\")].cluster.server}'`
	api_server_url.sub("https://", "")
end

def docker_account
	"docker-registry.#{dns_domain}"
end

def hostname(stack)
	if stack == @blue_stack
		"#{@app_name}-test.#{dns_domain}"
	else
		"#{@app_name}.#{dns_domain}"
	end
end

def get_version
	if File.exist?(@version_file)
		file = File.read(@version_file)
		require 'json'
		hash = JSON.parse(file)
		return hash["version"]
	end

	raise "Version file (#{@version_file} does not exist)"
end

def smoke_test(hostname)
	health_url = "#{hostname}/version"
	puts "Calling url: #{health_url}"

	status = `curl -s -o /dev/null --max-time 3 -w %{http_code} #{health_url}`.strip
	puts "\tStatus code: #{status}"
	return status == "200"
end

def smoke_test_timer(hostname, timeout)
	 # Keep trying for 'n' secs for smoke tests to pass. This gives time to nginx to pickup ingress changes
  begin
    require 'timeout'
    Timeout::timeout(timeout) do
	    while !smoke_test(hostname)
	      sleep 2
	    end
    	return true
  	end
  rescue Timeout::Error
		puts "Timed out waiting waiting for smoke test to pass"
		return false
	end
end

def wait_for_pods(filter)
	begin
		require 'timeout'
		require 'json'

		Timeout::timeout(@pod_timeout_secs) do
			begin
				sleep 2
				json_data = `kubectl get pods #{filter} -o json`
				data = JSON.parse(json_data)
				expected_count = data["items"].count
				running_count = data["items"].select{ |i| i["status"]["conditions"] != nil && i["status"]["conditions"].select{ |c| c["type"] == "Ready" && c["status"] == "True"}.count > 0}.count

				puts "Waiting for pods to be ready (#{running_count}/#{expected_count})..."
			end while running_count != expected_count
		end
	rescue Timeout::Error
    raise "Timed out waiting for new pods to be ready"
  end
end

def delete_stack(filter)
	blue_stack_exists = (`kubectl get ingress #{filter} -o name` != "")

	# If blue stack already exist, delete it
	if blue_stack_exists
		puts "Deleting older blue stack..."
		sh "kubectl delete ingress,svc,rc #{filter}"
		# Change stack label for pods as it takes time to terminate
		sh "kubectl label pods #{filter} --overwrite stack=terminate"
	end
end

def get_template_vars
	green_filter = "-l app=#{@app_name},stack=#{@green_stack}"
	green_stack_exists = (`kubectl get ingress #{green_filter} -o name` != "")

	# If new service(no green stack), then don't append hex to name
	if green_stack_exists
		require "securerandom"
		deploy_hex = SecureRandom.hex(4)
		name = "#{@app_name}-#{deploy_hex}"
	else
		name = "#{@app_name}"
	end

	blue_host = hostname(@blue_stack)
	vars = {
		:name => name,
		:app_name => @app_name,
		:stack => @blue_stack,
		:deploy_hex => deploy_hex,
		:account => docker_account,
		:image => @app_name,
		:version => get_version,
		:service_port => @service_port,
		:container_port => @container_port,
		:hostname => blue_host
	}

	return vars
end

task :test do
	vars = get_template_vars
	puts vars["name"]
	puts "Name: #{vars[:name]}"
end

task :deploy_blue do
	begin
		blue_filter = "-l app=#{@app_name},stack=#{@blue_stack}"
		delete_stack(blue_filter)

		vars = get_template_vars
		render_templates(vars)

		puts "\nCreating blue stack with name #{vars[:name]}..."
		# Create new blue stack
		sh "kubectl create -f #{@deploy_directory}/svc.yaml"
		sh "kubectl create -f #{@deploy_directory}/rc.yaml"
		sh "kubectl create -f #{@deploy_directory}/ingress.yaml"

		wait_for_pods(blue_filter)
	rescue Exception => e
		raise "Unexpected error while deploying blue stack. ERROR: #{e.message}"
	end
end

task :smoke_test_blue do
	blue_host = hostname(@blue_stack)
	if !smoke_test(blue_host)
		raise "Smoke test failed on blue stack. Deployment is stopped"
	end

	puts "Blue stack is healthy.\n"
end

task :promote_blue_to_green do
	begin
		blue_filter = "-l app=#{@app_name},stack=#{@blue_stack}"
		green_filter = "-l app=#{@app_name},stack=#{@green_stack}"
		green_stack_exists = (`kubectl get ingress #{green_filter} -o name` != "")
		green_host = hostname(@green_stack)

		puts "\nDirecting production traffic to blue stack..."
		require 'json'
		if green_stack_exists
			green_svc = `kubectl get svc #{green_filter} -o json`
			blue_svc = `kubectl get svc #{blue_filter} -o json`

			green_svc_hash = JSON.parse(green_svc)
			blue_svc_hash = JSON.parse(blue_svc)

			if blue_svc_hash["items"].empty? || green_svc_hash["items"].empty?
				raise "Cannot read blue or green service."
			end

			blue_svc_name = blue_svc_hash["items"].first["metadata"]["name"]
			# Update green svc to point to blue pods (use blue svc definition expect name,clusterIP)
			blue_svc_hash["items"].first["metadata"]["name"] = green_svc_hash["items"].first["metadata"]["name"]
			blue_svc_hash["items"].first["spec"]["clusterIP"] = green_svc_hash["items"].first["spec"]["clusterIP"]
			blue_svc_hash["items"].first["metadata"].delete("uid")
			blue_svc_hash["items"].first["metadata"].delete("resourceVersion")
			blue_svc_hash["items"].first["spec"]["ports"].map{ |p| p.delete("nodePort") }
			sh "echo '#{blue_svc_hash["items"].first.to_json}' | kubectl replace -f -"

			# delete green rc as traffic is now flowing to blue rc
			sh "kubectl delete rc #{green_filter}"
			# delete blue ingress,svc
			sh "kubectl delete ingress,svc #{blue_svc_name}"

			if !smoke_test(green_host)
				raise "Smoke test failed after promoting blue stack. Deployment failed."
			end
		else
			blue_ingress = `kubectl get ingress #{blue_filter} -o json`
			blue_ingress_hash = JSON.parse(blue_ingress)

			if blue_ingress_hash["items"].empty?
				raise "Cannot read blue ingress."
			end

			name = blue_ingress_hash["items"].first["metadata"]["name"]
			# update hostname(ingress) for blue stack to match hostname of green stack
			blue_ingress_hash["items"].first["spec"]["rules"].first["host"] = green_host
			sh "kubectl patch ingress #{name} -p '{ \"spec\": #{blue_ingress_hash["items"].first["spec"].to_json} }'"

			# Smoke test with timer to allow nginx to pickup ingress hcnage
			if !smoke_test_timer(green_host, 20)
				raise "Smoke test failed after promoting blue stack. Deployment failed."
			end
		end

		# Update stack label from blue to green
		sh "kubectl label ingress,svc,rc,pods #{blue_filter} --overwrite stack=#{@green_stack}"
	rescue Exception => e
		raise "Unexpected error when promoting blue stack to green. ERROR: #{e.message}"
	end
end

task :deploy => [:deploy_blue, :smoke_test_blue, :promote_blue_to_green]

task :clean do
	FileUtils.rm_rf(@build_directory)
	FileUtils.rm_rf(@deploy_directory)
end

class Namespace
  def initialize(hash)
    hash.each do |key, value|
      singleton_class.send(:define_method, key) { value }
    end
  end

  def get_binding
    binding
  end
end
