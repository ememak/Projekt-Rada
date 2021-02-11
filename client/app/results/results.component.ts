import { Component } from '@angular/core';
import { ActivatedRoute,
         Router } from '@angular/router';

import { grpc } from '@improbable-eng/grpc-web';
import { Query } from "Projekt_Rada/query/query_pb_service";
import { PollSchema, PollSummary, SummaryRequest } from "Projekt_Rada/query/query_pb";
import { host } from '../host';

@Component({
  selector: 'results',
  templateUrl: './results.component.html',
})
export class ResultsComponent {
  summary: PollSummary.AsObject;

  input: any[] = [];

  pollid: number;
  inpid: number;

  constructor (private route: ActivatedRoute, private router: Router) {
    this.pollid = parseInt(this.route.snapshot.paramMap.get('pollid'));
    if(!isNaN(this.pollid)){
    let request: SummaryRequest = new SummaryRequest();
    request.setPollid(this.pollid);
    grpc.unary(Query.GetSummary, {
      request: request,
      host: host,
      onEnd: res => {
        const { status, statusMessage, headers, message, trailers } = res;
        console.log("pollInit.onEnd.status", status, statusMessage);
        console.log("pollInit.onEnd.headers", headers);
        if (status === grpc.Code.OK && message) {
          console.log("pollInit.onEnd.message", message.toObject());
          this.summary = (<PollSummary>message).toObject()
          for (let qa of this.summary.schema.questionsList) {
            let r = [];
            for(let [i, opt] of qa.optionsList.entries()) {
              r.push([opt, parseInt(qa.answersList[i])]);
            }
            this.input.push(r)
          }
        }
        console.log("pollInit.onEnd.trailers", trailers);
      }
    });}
  }

  getPollid() {
    this.router.navigate(['/results', this.inpid]);
  }

  trackOption(index: number, option: string) {
    return index;
  }

  get diagnostic() { return JSON.stringify(this.summary.schema); }
  
  onSubmit() {}
}
