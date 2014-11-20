FROM ubuntu:latest

MAINTAINER Abhi Yerra <abhi@berkeley.edu>

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y nginx openssh-server git-core openssh-client curl
RUN apt-get install -y build-essential
RUN apt-get install -y openssl libreadline6 libreadline6-dev zlib1g zlib1g-dev libssl-dev libyaml-dev libsqlite3-dev sqlite3 libxml2-dev libxslt-dev autoconf libc6-dev ncurses-dev automake libtool bison pkg-config
RUN apt-get install -y postgresql-client

# RUN useradd -ms /bin/bash deployer
# USER deployer

RUN gpg --keyserver hkp://keys.gnupg.net --recv-keys D39DC0E3
RUN \curl -L https://get.rvm.io | bash -s stable
RUN /bin/bash -l -c "rvm install 2.1"
RUN /bin/bash -l -c "echo 'gem: --no-ri --no-rdoc' > ~/.gemrc"
RUN /bin/bash -l -c "gem install bundler --no-ri --no-rdoc"
RUN apt-get install -y libpq5 libpq-dev

RUN mkdir -p /treely
ADD . /treely
WORKDIR /treely

RUN /bin/bash -l -c "bundle"

ENV RAILS_ENV production
CMD /bin/bash -l -c "RAILS_ENV=production bundle exec rake assets:precompile"
CMD /bin/bash -l -c "bundle exec rails s -e production -p 3000"