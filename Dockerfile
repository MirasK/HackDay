###################################################################################################
#												                                                  #
#                                   Miron-developer & Mirask                                      #
#				                          Hackday2020                       					  #
#												                                                  #
###################################################################################################

FROM golang:onbuild

WORKDIR /go/src/app
COPY . .

LABEL description="This is the hackday2020 project." \
    authors="Miron-developer, MirasK" \
    contacts="https://github.com/miron-developer, https://github.com/mirasK"

CMD ["app"]

EXPOSE 8080