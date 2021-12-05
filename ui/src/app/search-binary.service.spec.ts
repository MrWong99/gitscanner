import { TestBed } from '@angular/core/testing';

import { SearchBinaryService } from './search-binary.service';

describe('SearchBinaryService', () => {
  let service: SearchBinaryService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(SearchBinaryService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
