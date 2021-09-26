curl -XPOST http://localhost:8080/sms-query --data "text=help" --data "msisdn=33631396906"

curl -XPOST http://localhost:8080/sms-query --data "text=search 'orange'" --data "msisdn=33631396906"

curl -XPOST http://localhost:8080/sms-query --data "text=translate en fr 'plane and car'" --data "msisdn=33631396906"

curl -XPOST http://localhost:8080/sms-query --data "text=news" --data "msisdn=33631396906"
curl -XPOST http://localhost:8080/sms-query --data "text=lemonde" --data "msisdn=33631396906"
curl -XPOST http://localhost:8080/sms-query --data "text=lemonde politique" --data "msisdn=33631396906"

curl -XPOST http://localhost:8080/sms-query --data "text=weather" --data "msisdn=33631396906"
curl -XPOST http://localhost:8080/sms-query --data "text=weather hour" --data "msisdn=33631396906"
curl -XPOST http://localhost:8080/sms-query --data "text=weather bordeaux" --data "msisdn=33631396906"
curl -XPOST http://localhost:8080/sms-query --data "text=weather hour bordeaux" --data "msisdn=33631396906"

curl -XPOST http://localhost:8080/sms-query --data "text=rain" --data "msisdn=33631396906"
curl -XPOST http://localhost:8080/sms-query --data "text=rain bordeaux" --data "msisdn=33631396906"

curl -XPOST http://localhost:8080/sms-query --data "text=bicloo 'gare maritime'" --data "msisdn=33631396906"
curl -XPOST http://localhost:8080/sms-query --data "text=velib soljenitsyne" --data "msisdn=33631396906"

curl -XPOST http://localhost:8080/sms-query --data "text=tan t1 lauriers beaujoire" --data "msisdn=33631396906"

curl -XPOST http://localhost:8080/sms-query --data "text=stif b61 'Eglise de Pantin' 'Gare d'Austerlitz'" --data "msisdn=33631396906"

# aliases
curl -XPOST http://localhost:8080/sms-query --data "text=bicloo" --data "msisdn=33631396906"
curl -XPOST http://localhost:8080/sms-query --data "text=downtown" --data "msisdn=33631396906"
curl -XPOST http://localhost:8080/sms-query --data "text=home" --data "msisdn=33631396906"
curl -XPOST http://localhost:8080/sms-query --data "text=travail" --data "msisdn=33666559968"
curl -XPOST http://localhost:8080/sms-query --data "text=maison" --data "msisdn=33666559968"
