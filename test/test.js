import http from "k6/http";
import { check } from "k6";
import { Rate } from "k6/metrics";

export let errorRate = new Rate("errors");

export default function() {
  var url = __ENV.MY_HOSTNAME;
  var params = {  };
  check(http.get(url, params), {
    "status is 200": (r) => r.status == 200
  }) || errorRate.add(1);
};