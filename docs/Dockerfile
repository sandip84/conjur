FROM ruby:3.0.6-slim-bullseye
RUN gem update --system
RUN gem install bundler jekyll
RUN mkdir /src
COPY ./Gemfile** /src/
WORKDIR /src
RUN bundle install
