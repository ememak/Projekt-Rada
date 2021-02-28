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

  // Each entry in this array is an input for graph regarding question with matching index.
  // Graph input is array of arrays consisting of at pair [key, value].
  graphsInput: any[] = [];

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
          if (status === grpc.Code.OK && message) {
            this.summary = (<PollSummary>message).toObject()
            this.getGraphsInput()
          }
        }
      });
    }
  }

  getPollid() {
    this.router.navigate(['/results', this.inpid]);
  }

  trackOption(index: number, option: string) {
    return index;
  }

  getGraphsInput() {
    for (let qa of this.summary.schema.questionsList) {
      let graphInp = []; 
      for(let [i, opt] of qa.optionsList.entries()) {
        graphInp.push([opt, parseInt(qa.answersList[i])]);
      }
      this.graphsInput.push(graphInp)
    }
    console.log(this.summary)
  }

  get diagnostic() { return JSON.stringify(this.summary); }

  onSubmit() {}
}
