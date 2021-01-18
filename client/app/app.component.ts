import { Component } from '@angular/core';

import { grpc } from '@improbable-eng/grpc-web';
import { Query } from "Projekt_Rada/query/query_pb_service";
import { PollSchema } from "Projekt_Rada/query/query_pb";

const host = "http://localhost:8080";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: [ './app.component.css' ]
})
export class AppComponent  {
  pollInit() {
    const schema = new PollSchema();
    console.log(schema)
    grpc.unary(Query.PollInit, {
      request: schema,
      host: host,
      metadata: new grpc.Metadata({"Access-Control-Allow-Origin": "*"}),
      onEnd: res => {
        const { status, statusMessage, headers, message, trailers } = res;
        console.log("pollInit.onEnd.status", status, statusMessage);
        console.log("pollInit.onEnd.headers", headers);
        if (status === grpc.Code.OK && message) {
          console.log("pollInit.onEnd.message", message.toObject());
        }
        console.log("pollInit.onEnd.trailers", trailers);
      }
    });
  }
}

/*
Copyright Google LLC. All Rights Reserved.
Use of this source code is governed by an MIT-style license that
can be found in the LICENSE file at https://angular.io/license
*/
