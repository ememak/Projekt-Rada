import { Component } from '@angular/core';

import { grpc } from '@improbable-eng/grpc-web';
import { Query } from "Projekt_Rada/query/query_pb_service";
import { PollSchema } from "Projekt_Rada/query/query_pb";

const host = "http://localhost:12345";

@Component({
  selector: 'poll-init',
  templateUrl: './pollinit.component.html',
})
export class PollInitComponent {
  questionsList: PollSchema.QA.AsObject[]  = [{
      question: "",
      type: PollSchema.QuestionType.OPEN,
      answer: "",
    },
  ];

  constructor () {}

  addQuestion() {
    this.questionsList.push({
      question: "",
      type: PollSchema.QuestionType.OPEN,
      answer: "",
    });
  }

  get diagnostic() { return JSON.stringify(this.questionsList); }

  onSubmit() {}
}


/*
Copyright Google LLC. All Rights Reserved.
Use of this source code is governed by an MIT-style license that
can be found in the LICENSE file at https://angular.io/license
*/
