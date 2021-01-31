import { Component } from '@angular/core';

import { grpc } from '@improbable-eng/grpc-web';
import { Query } from "Projekt_Rada/query/query_pb_service";
import { PollSchema } from "Projekt_Rada/query/query_pb";

const host = "http://localhost:12345";

@Component({
  selector: 'app-pollinit',
  templateUrl: './pollinit.component.html',
})
export class PollInitComponent {
  input: { question: string, qtype: number }[] = [
    { "question": "", "qtype": 0 },
  ];
  
  ngOnInit() {
  }

  newQuestion() {
    this.input.push({"question": "", "qtype": 0});
  }

  sendPoll() {
    const schema= new PollSchema();
    for (let inp of this.input){
      const QA = new PollSchema.QA();
      QA.setQuestion(inp.question);
      QA.setType(inp.qtype as 0 | 1 | 2);
      schema.addQuestions(QA);
    }
    grpc.unary(Query.PollInit, {
      request: schema,
      host: host,
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
