import { Component, OnInit, Input } from '@angular/core';
import { FileData, Error } from '../search-binary.service';

@Component({
  selector: 'app-scan-overview',
  templateUrl: './scan-overview.component.html',
  styleUrls: ['./scan-overview.component.scss']
})

export class ScanOverviewComponent implements OnInit {

  @Input() data: any;
  @Input() checkName: any;
  errors: Error[] = [{repository: 'gus', error: 'you suck'}];

  constructor() { }

  ngOnInit(): void {    
  }

  isBinariesSearch() {
    return this.checkName.includes('SearchBinaries');
  }

  isUnicodeSearch() {
    return this.checkName.includes('SearchUnicode');
  }

  isCommitSearch() {
    return this.checkName.includes('CheckCommitAuthor');
  }

  getCheckName(checkNamePath: string): string {
    checkNamePath.slice(checkNamePath.lastIndexOf('.') + 1);
    return checkNamePath;
  }

  listErrors() {    
    this.data.array.forEach((entry: FileData) => {
      if (entry.error && entry.repository && entry.error !== '') {
        this.errors.push({repository: entry.repository, error: entry.error})
      }
    });    
  }

}
